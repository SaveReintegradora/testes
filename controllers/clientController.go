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
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

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
		fmt.Println("[ERRO] Arquivo não enviado:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Arquivo não enviado"})
		return
	}
	fmt.Println("Arquivo recebido:", file.Filename, file.Size)

	// Verifica extensão
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xls" && ext != ".xlsx" {
		fmt.Println("[ERRO] Extensão inválida:", ext)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "O arquivo deve ser .xls ou .xlsx"})
		return
	}

	tempPath := "/tmp/" + uuid.New().String() + ext
	if err := ctx.SaveUploadedFile(file, tempPath); err != nil {
		fmt.Println("[ERRO] Erro ao salvar arquivo temporário:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao salvar arquivo temporário"})
		return
	}
	defer os.Remove(tempPath)

	var rows [][]string
	if ext == ".xlsx" {
		xl, err := excelize.OpenFile(tempPath)
		if err != nil {
			fmt.Println("[ERRO] Erro ao abrir arquivo Excel:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao abrir arquivo Excel"})
			return
		}
		sheet := xl.GetSheetName(0)
		rows, err = xl.GetRows(sheet)
		if err != nil {
			fmt.Println("[ERRO] Erro ao ler linhas do Excel:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao ler linhas do Excel"})
			return
		}
	} else {
		xlsFile, err := xls.Open(tempPath, "utf-8")
		if err != nil {
			fmt.Println("[ERRO] Erro ao abrir arquivo XLS:", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erro ao abrir arquivo XLS"})
			return
		}
		sheet := xlsFile.GetSheet(0)
		if sheet == nil {
			fmt.Println("[ERRO] Não foi possível ler a primeira planilha do XLS")
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
		fmt.Println("[ERRO] Arquivo Excel vazio")
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
		} else if normCol == "email" || normCol == "email" {
			colMap["email"] = idx
		} else if normCol == "telefone" {
			colMap["telefone"] = idx
		} else if normCol == "endereco" {
			colMap["endereco"] = idx
		}
	}
	fmt.Println("Mapeamento de colunas:", colMap)
	// Verifica se todos os campos obrigatórios existem
	required := []string{"nome", "email", "telefone", "endereco"}
	for _, req := range required {
		if _, ok := colMap[req]; !ok {
			fmt.Println("[ERRO] Cabeçalho faltando campo obrigatório:", req)
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
		if colMap["nome"] < len(row) {
			name = row[colMap["nome"]]
		}
		if colMap["email"] < len(row) {
			email = row[colMap["email"]]
		}
		if colMap["telefone"] < len(row) {
			phone = row[colMap["telefone"]]
		}
		if colMap["endereco"] < len(row) {
			address = row[colMap["endereco"]]
		}
		// Tenta pegar CNPJ se existir na planilha
		for idx, col := range header {
			if normalizeHeader(col) == "cnpj" && idx < len(row) {
				cnpj = row[idx]
			}
		}
		fmt.Printf("[DEBUG] Linha %d: nome='%s', cnpj='%s', email='%s', telefone='%s', endereco='%s'\n", i, name, cnpj, email, phone, address)
		// Validação: não cadastrar cliente duplicado
		var exists bool
		if cnpj != "" {
			exists, err = c.repo.ExistsByNameAndCNPJ(name, cnpj)
		} else {
			exists, err = c.repo.ExistsByNameAndEmail(name, email)
		}
		if err != nil {
			if !strings.Contains(err.Error(), "record not found") {
				fmt.Printf("[ERRO] Falha ao verificar duplicidade (linha %d): %v\n", i, err)
				dbErrors++
			}
			continue
		}
		if exists {
			fmt.Printf("[INFO] Cliente duplicado ignorado: nome='%s', cnpj='%s'\n", name, cnpj)
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
		if err := c.repo.Create(&client); err == nil {
			count++
		} else {
			fmt.Printf("[ERRO] Falha ao inserir cliente (linha %d): %v\n", i, err)
			dbErrors++
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"clientes importados": count,
		"ignorados por duplicidade": ignored,
		"erros de banco": dbErrors,
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
