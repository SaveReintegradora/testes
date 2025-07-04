import React, { useState } from "react";
import 'bootstrap/dist/css/bootstrap.min.css';

function App() {
  const [file, setFile] = useState(null);
  const [uploadStatus, setUploadStatus] = useState("");
  const [uploadError, setUploadError] = useState("");
  const [downloadUrl, setDownloadUrl] = useState("");
  const [exportStatus, setExportStatus] = useState("");
  const [exportError, setExportError] = useState("");

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
    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch("http://localhost:5000/clients/upload", {
        method: "POST",
        body: formData,
      });
      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || "Erro ao enviar arquivo");
      }
      setUploadStatus("Arquivo enviado com sucesso!");
    } catch (err) {
      setUploadError(err.message);
      setUploadStatus("");
    }
  };

  const exportClients = async () => {
    setExportStatus("Exportando...");
    setExportError("");
    setDownloadUrl("");
    try {
      const response = await fetch("http://localhost:5000/clients/export");
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
                  />
                  <button className="btn btn-primary" type="submit">
                    Enviar
                  </button>
                </div>
              </form>
              {uploadStatus && <div className="alert alert-success">{uploadStatus}</div>}
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

export default App;
