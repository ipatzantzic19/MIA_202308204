import { useState, useRef } from 'react';
import { Upload, Play, Trash2, Terminal, FileText, HardDrive } from 'lucide-react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import LoginPage from './pages/Login';
import styles from './styles/App.module.css';
import PanelEditor from './components/PanelEditor';
import PanelConsola from './components/PanelConsola';
import InfoCards from './components/InfoCards';

function App() {
  const [inputCommands, setInputCommands] = useState('');
  const [outputCommands, setOutputCommands] = useState('');
  const [isExecuting, setIsExecuting] = useState(false);
  const [fileName, setFileName] = useState('');
  const fileInputRef = useRef(null);

  // Función para cargar archivo de script
  const handleFileUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      // Verificar que sea un archivo .smia
      if (!file.name.endsWith('.smia')) {
        setOutputCommands(prev =>
          prev + `\n[ERROR] El archivo debe tener extensión .smia\n`
        );
        return;
      }

      const reader = new FileReader();
      reader.onload = (e) => {
        const content = e.target.result;
        setInputCommands(content);
        setFileName(file.name);
        setOutputCommands(prev =>
          prev + `\n[INFO] Archivo "${file.name}" cargado exitosamente\n`
        );
      };
      reader.onerror = () => {
        setOutputCommands(prev =>
          prev + `\n[ERROR] Error al leer el archivo\n`
        );
      };
      reader.readAsText(file);
    }
  };

  // Función para ejecutar comandos
  const handleExecute = async () => {
    if (!inputCommands.trim()) {
      setOutputCommands(prev =>
        prev + `\n[ERROR] No hay comandos para ejecutar\n`
      );
      return;
    }

    setIsExecuting(true);
    setOutputCommands(prev =>
      prev + `\n${'='.repeat(60)}\n[INICIO DE EJECUCIÓN]\n${'='.repeat(60)}\n`
    );

    try {
      // Aquí se enviará al backend cuando esté listo
      const response = await fetch('http://localhost:8080/api/execute', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ commands: inputCommands }),
      });

      if (response.ok) {
        const data = await response.json();
        setOutputCommands(prev => prev + data.output + '\n');
      } else {
        throw new Error('Error en la comunicación con el servidor');
      }
    } catch (error) {
      // Por ahora, simular el procesamiento de comandos
      const lines = inputCommands.split('\n');
      let output = '';

      lines.forEach((line, index) => {
        const trimmedLine = line.trim();

        // Ignorar líneas vacías
        if (!trimmedLine) {
          output += '\n';
          return;
        }

        // Mostrar comentarios
        if (trimmedLine.startsWith('#')) {
          output += `${trimmedLine}\n`;
          return;
        }

        // Simular procesamiento de comandos (será reemplazado por el backend)
        output += `[LÍNEA ${index + 1}] Procesando: ${trimmedLine}\n`;
        output += `[INFO] Comando enviado al servidor (Backend no conectado)\n\n`;
      });

      setOutputCommands(prev => prev + output);
    }

    setOutputCommands(prev =>
      prev + `${'='.repeat(60)}\n[FIN DE EJECUCIÓN]\n${'='.repeat(60)}\n`
    );
    setIsExecuting(false);
  };

  // Función para limpiar el área de salida
  const handleClearOutput = () => {
    setOutputCommands('');
  };

  // Función para limpiar el área de entrada
  const handleClearInput = () => {
    setInputCommands('');
    setFileName('');
  };

  return (
    <Router>
      <Routes>
        <Route
          path="/login"
          element={<LoginPage />}
        />
        <Route
          path="/"
          element={

            <div className={styles.appContainer}>

              {/* Header */}
              <header className={styles.header}>
                <div className={styles.headerContainer}>
                  <div className={styles.headerContent}>
                    <div className={styles.logoContainer}>
                      <HardDrive className={styles.logoIcon} />
                      <h1 className={styles.title}>GoDisk</h1>
                    </div>
                    <div className={styles.userInfo}>
                      <Link to="/login" className={styles.userButton}>Iniciar Sesión</Link>
                    </div>
                  </div>
                </div>
              </header>

              {/* Main Content */}
              <main className={styles.mainContent}>
                <div className={styles.gridContainer}>
                  <PanelEditor
                    inputCommands={inputCommands}
                    setInputCommands={setInputCommands}
                    handleFileUpload={handleFileUpload}
                    handleExecute={handleExecute}
                    handleClearInput={handleClearInput}
                    isExecuting={isExecuting}
                    fileInputRef={fileInputRef}
                    fileName={fileName}
                  />

                  <PanelConsola
                    outputCommands={outputCommands}
                    handleClearOutput={handleClearOutput}
                    styles={styles}
                  />

                </div>

                {/* Info Cards */}
                <div className={styles.infoCardsContainer}>
                  <InfoCards Card={{
                    title: "Comandos Disponibles",
                    content: "mkdisk, rmdisk, fdisk, mount, mounted, mkfs, login, logout, mkgrp, mkusr, mkfile, mkdir, cat, rep"
                  }} />

                  <InfoCards Card={{
                    title: "Formato de Script",
                    content: "Los scripts deben tener extensión .smia. Use # para comentarios. Los parámetros se separan por espacios."
                  }} />

                  <InfoCards Card={{
                    title: "Estado del Sistema",
                    content: isExecuting ? (
                      <span className={styles.statusExecuting}>⚡ Ejecutando comandos...</span>
                    ) : (
                      <span className={styles.statusReady}>✓ Listo para ejecutar</span>
                    )
                  }} />

                </div>
              </main>
            </div>
          } />
      </Routes>
    </Router>
  );
}

export default App;
