import React, { useState, useRef } from "react";
import 'bootstrap/dist/css/bootstrap.min.css';

// Exemplo de função para consumir CRUD de clientes com API Key
// Pode ser movido para src/api.js
export async function getClients(API_URL, API_KEY) {
  const response = await fetch(`${API_URL}/clients`, {
    headers: { 'x-api-key': API_KEY }
  });
  if (!response.ok) throw new Error('Erro ao buscar clientes');
  return response.json();
}

function App() {
  const [file, setFile] = useState(null);
  const [uploadStatus, setUploadStatus] = useState("");
  const [uploadError, setUploadError] = useState("");
  const [downloadUrl, setDownloadUrl] = useState("");
  const [exportStatus, setExportStatus] = useState("");
  const [exportError, setExportError] = useState("");
  const [duplicados, setDuplicados] = useState(0);
  const fileInputRef = useRef();

  // Use variável de ambiente para a URL da API
  const API_URL = process.env.REACT_APP_API_URL || "http://api:5000";
  const API_KEY = process.env.REACT_APP_API_KEY || "SUA_API_KEY_AQUI"; // Defina no .env

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
    setUploadStatus("");
    setUploadError("");
  };

  const handleUpload = async (e) => {
    e.preventDefault();
    if (!file) {
      setUploadError("Selecione um arquivo Excel (.xls ou .xlsx)");
      return;
    }
    setUploadStatus("Enviando...");
    setUploadError("");
    setDuplicados(0);
    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch(`${API_URL}/clients/upload`, {
        method: "POST",
        body: formData,
        headers: {
          'x-api-key': API_KEY
        },
      });
      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.error || "Erro ao enviar arquivo");
      }
      let msg = `Clientes importados: ${data["clientes importados"]}`;
      if (data["ignorados por duplicidade"]) {
        setDuplicados(data["ignorados por duplicidade"]);
        msg += ` | Duplicados ignorados: ${data["ignorados por duplicidade"]}`;
      } else {
        setDuplicados(0);
      }
      if (data["erros de banco"]) {
        msg += ` | Erros de banco: ${data["erros de banco"]}`;
      }
      setUploadStatus(msg);
      setFile(null); // Limpa o estado do arquivo
      if (fileInputRef.current) fileInputRef.current.value = ""; // Limpa input
    } catch (err) {
      setUploadError(err.message);
      setUploadStatus("");
      setDuplicados(0);
    }
  };

  const exportClients = async () => {
    setExportStatus("Exportando...");
    setExportError("");
    setDownloadUrl("");
    try {
      const response = await fetch(`${API_URL}/clients/export`, {
        headers: {
          'x-api-key': API_KEY
        }
      });
      if (!response.ok) throw new Error("Erro ao exportar clientes");
      const data = await response.json();
      setDownloadUrl(data.download_url);
      setExportStatus("Arquivo pronto para download!");
    } catch (err) {
      setExportError(err.message || "Erro desconhecido");
      setExportStatus("");
    }
  };

  return (
    <div className="container py-5">
      <div className="text-center mb-5">
        <h1 className="display-4 fw-bold" style={{ color: "#0d6efd" }}>Save Reintegradora</h1>
        <p className="lead">Gestão de Clientes - Importação e Exportação</p>
      </div>

      <div className="row justify-content-center">
        <div className="col-md-6">
          <div className="card shadow-sm mb-4">
            <div className="card-body">
              <h3 className="card-title mb-3">Importar Clientes (Excel)</h3>
              <form onSubmit={handleUpload}>
                <div className="input-group mb-3">
                  <input
                    type="file"
                    className="form-control"
                    accept=".xls,.xlsx"
                    onChange={handleFileChange}
                    ref={fileInputRef}
                  />
                  <button className="btn btn-primary" type="submit">
                    Enviar
                  </button>
                </div>
              </form>
              {uploadStatus && <div className="alert alert-success">{uploadStatus}</div>}
              {duplicados > 0 && (
                <div className="alert alert-warning">{`Atenção: ${duplicados} cliente(s) já existiam e foram ignorados.`}</div>
              )}
              {uploadError && <div className="alert alert-danger">{uploadError}</div>}
            </div>
          </div>

          <div className="card shadow-sm">
            <div className="card-body">
              <h3 className="card-title mb-3">Exportar Clientes</h3>
              <button className="btn btn-success" onClick={exportClients}>
                Exportar para XLS
              </button>
              {exportStatus && <div className="alert alert-success mt-3">{exportStatus}</div>}
              {exportError && <div className="alert alert-danger mt-3">{exportError}</div>}
              {downloadUrl && (
                <div className="mt-3">
                  <a
                    href={downloadUrl}
                    className="btn btn-outline-primary"
                    target="_blank"
                    rel="noopener noreferrer"
                  >
                    Baixar arquivo XLS
                  </a>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

// Sugestão: crie um arquivo src/api.js para centralizar chamadas de API
// Exemplo:
// export async function uploadClients(file) { ... }
// export async function exportClients() { ... }
// Isso facilita testes e manutenção.

// Sugestão: crie testes com Jest/React Testing Library para upload/exportação e feedback visual.

export default App;
