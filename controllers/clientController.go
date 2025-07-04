package controllers

import (
	"fmt"
	"minha-api/models"
	"minha-api/repositories"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/extrame/xls"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

var validate = validator.New()

type ClientController struct {
	repo *repositories.ClientRepository
}

func NewClientController(repo *repositories.ClientRepository) *ClientController {
	return &ClientController{repo: repo}
}

// UploadClients godoc
// @Summary      Upload de clientes via arquivo Excel
// @Description  Recebe um arquivo .xls, lê os dados e cadastra clientes no banco
// @Tags         clients
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Arquivo de clientes (.xls)"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /clients/upload [post]
func (c *ClientController) UploadClients(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não enviado"})
		return
	}
	fmt.Println("Arquivo recebido:", file.Filename, file.Size)

	// Verifica extensão
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xls" && ext != ".xlsx" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "O arquivo deve ser .xls ou .xlsx"})
		return
	}

	tempPath := "/tmp/" + uuid.New().String() + ext
	if err := ctx.SaveUploadedFile(file, tempPath); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao salvar arquivo temporário"})
		return
	}
	defer os.Remove(tempPath)

	var rows [][]string
	if ext == ".xlsx" {
		xl, err := excelize.OpenFile(tempPath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao abrir arquivo Excel"})
			return
		}
		// Busca a primeira aba não vazia
		sheetName := ""
		for _, name := range xl.GetSheetList() {
			rowsTmp, _ := xl.GetRows(name)
			if len(rowsTmp) > 0 {
				sheetName = name
				rows = rowsTmp
				break
			}
		}
		if sheetName == "" || len(rows) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo Excel vazio ou sem dados válidos"})
			return
		}
	} else {
		xlsFile, err := xls.Open(tempPath, "utf-8")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao abrir arquivo XLS"})
			return
		}
		sheet := xlsFile.GetSheet(0)
		if sheet == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Não foi possível ler a primeira planilha do XLS"})
			return
		}
		for i := 0; i <= int(sheet.MaxRow); i++ {
			row := sheet.Row(i)
			var rowData []string
			for j := 0; j < row.LastCol(); j++ {
				rowData = append(rowData, row.Col(j))
			}
			rows = append(rows, rowData)
		}
	}

	if len(rows) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo Excel vazio"})
		return
	}

	header := rows[0]
	fmt.Println("Header detectado:", header)
	// Mapeamento flexível dos campos
	colMap := map[string]int{}
	for idx, col := range header {
		normCol := normalizeHeader(col)
		fmt.Printf("Coluna original: '%s' | Normalizada: '%s'\n", col, normCol)
		if normCol == "nome" {
			colMap["nome"] = idx
		} else if normCol == "email" {
			colMap["email"] = idx
		} else if normCol == "telefone" {
			colMap["telefone"] = idx
		} else if normCol == "endereco" {
			colMap["endereco"] = idx
		} else if normCol == "cnpj" {
			colMap["cnpj"] = idx
		} // ignora qualquer outra coluna (ex: id)
	}
	fmt.Println("Mapeamento de colunas:", colMap)
	// Verifica se todos os campos obrigatórios existem
	required := []string{"nome", "email", "telefone", "endereco"}
	for _, req := range required {
		if _, ok := colMap[req]; !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cabeçalho do arquivo deve conter as colunas: Nome, Email, Telefone, Endereço (em qualquer ordem)"})
			return
		}
	}

	var count, ignored, dbErrors int
	for i, row := range rows {
		if i == 0 {
			continue // pula header
		}
		// Pega os valores pelas posições mapeadas
		name := ""
		email := ""
		phone := ""
		address := ""
		cnpj := ""
		if idx, ok := colMap["nome"]; ok && idx < len(row) {
			name = row[idx]
		}
		if idx, ok := colMap["email"]; ok && idx < len(row) {
			email = row[idx]
		}
		if idx, ok := colMap["telefone"]; ok && idx < len(row) {
			phone = row[idx]
		}
		if idx, ok := colMap["endereco"]; ok && idx < len(row) {
			address = row[idx]
		}
		if idx, ok := colMap["cnpj"]; ok && idx < len(row) {
			cnpj = row[idx]
		}
		fmt.Printf("[IMPORT] Linha %d: nome='%s', email='%s', telefone='%s', endereco='%s', cnpj='%s'\n", i, name, email, phone, address, cnpj)
		// Normaliza campos para evitar duplicidade por diferença de maiúsculas/minúsculas/espacos
		name = strings.TrimSpace(strings.ToLower(name))
		email = strings.TrimSpace(strings.ToLower(email))
		cnpj = strings.TrimSpace(strings.ToLower(cnpj))
		// Validação: não cadastrar cliente duplicado
		var exists bool
		if cnpj != "" {
			exists, err = c.repo.ExistsByNameAndCNPJ(name, cnpj)
		} else {
			exists, err = c.repo.ExistsByNameAndEmail(name, email)
		}
		if err != nil {
			dbErrors++
			continue
		}
		if exists {
			ignored++
			continue
		}
		client := models.Client{
			ID:      uuid.New().String(),
			Name:    name,
			Email:   email,
			Phone:   phone,
			Address: address,
			CNPJ:    cnpj,
		}
		// Validação robusta dos campos
		err := validate.StructPartial(client, "Name", "Email", "Phone", "Address")
		if err != nil {
			dbErrors++
			continue
		}
		if err := c.repo.Create(&client); err == nil {
			count++
		} else {
			errMsg := strings.ToLower(err.Error())
			if strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "unique") {
				ignored++
			} else {
				dbErrors++
			}
		}
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"clientes importados":       count,
		"ignorados por duplicidade": ignored,
		"erros de banco":            dbErrors,
	})
}

func normalizeHeader(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// Remove acentos e caracteres especiais
	norm := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			norm = append(norm, unicode.ToLower(r))
		}
	}
	return string(norm)
}
