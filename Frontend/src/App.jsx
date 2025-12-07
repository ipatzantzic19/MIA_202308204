import { useState, useRef } from 'react';
import { Upload, Play, Trash2, Terminal, FileText, HardDrive } from 'lucide-react';
import styles from './App.module.css';

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
    
    <div className={styles.appContainer}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContainer}>
          <div className={styles.headerContent}>
            <div className={styles.logoContainer}>
              <HardDrive className={styles.logoIcon} />
              <div>
                <h1 className={styles.title}>GoDisk</h1>
                <p className={styles.subtitle}>Sistema de Archivos EXT2</p>
              </div>
            </div>
            <div className={styles.userInfo}>
              <div className={styles.userInfoText}>
                <p className={styles.userInfoPrimary}>Universidad San Carlos</p>
                <p className={styles.userInfoSecondary}>MIA - Proyecto 1</p>
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className={styles.mainContent}>
        <div className={styles.gridContainer}>
          {/* Input Section */}
          <div className={styles.card}>
            <div className={styles.cardHeader}>
              <div className={styles.cardHeaderContent}>
                <div className={styles.cardTitleContainer}>
                  <Terminal className={`${styles.cardIcon} text-green-400`} />
                  <h2 className={styles.cardTitle}>Entrada de Comandos</h2>
                </div>
                {fileName && (
                  <span className={styles.fileName}>
                    <FileText className={styles.fileNameIcon} />
                    <span>{fileName}</span>
                  </span>
                )}
              </div>
            </div>

            <div className={styles.cardBody}>
              <textarea
                value={inputCommands}
                onChange={(e) => setInputCommands(e.target.value)}
                placeholder="# Ingrese sus comandos aquí o cargue un archivo .smia
# Ejemplo:
mkdisk -size=3000 -unit=K
fdisk -size=300 -diskName=VDIC-A.mia -name=Particion1"
                className={styles.textarea}
                spellCheck="false"
              />

              <div className={styles.buttonGroup}>
                <input
                  type="file"
                  ref={fileInputRef}
                  onChange={handleFileUpload}
                  accept=".smia"
                  className={styles.hiddenInput}
                />
                <button
                  onClick={() => fileInputRef.current?.click()}
                  className={`${styles.button} ${styles.buttonPrimary}`}
                >
                  <Upload className={styles.buttonIcon} />
                  <span>Cargar Script</span>
                </button>

                <button
                  onClick={handleExecute}
                  disabled={isExecuting || !inputCommands.trim()}
                  className={`${styles.button} ${styles.buttonSuccess}`}
                >
                  <Play className={styles.buttonIcon} />
                  <span>{isExecuting ? 'Ejecutando...' : 'Ejecutar'}</span>
                </button>

                <button
                  onClick={handleClearInput}
                  className={styles.buttonSecondary}
                  title="Limpiar entrada"
                >
                  <Trash2 className={styles.buttonIcon} />
                </button>
              </div>
            </div>
          </div>

          {/* Output Section */}
          <div className={styles.card}>
            <div className={styles.cardHeader}>
              <div className={styles.cardHeaderContent}>
                <div className={styles.cardTitleContainer}>
                  <Terminal className={`${styles.cardIcon} text-yellow-400`} />
                  <h2 className={styles.cardTitle}>Salida de Comandos</h2>
                </div>
                <button
                  onClick={handleClearOutput}
                  className={styles.clearButton}
                >
                  <Trash2 className={styles.buttonIcon} />
                  <span>Limpiar</span>
                </button>
              </div>
            </div>

            <div className={styles.cardBody}>
              <div className={styles.outputContainer}>
                <pre className={styles.preformatted}>
                  {outputCommands || '# Esperando ejecución de comandos...\n# Los resultados aparecerán aquí'}
                </pre>
              </div>
            </div>
          </div>
        </div>

        {/* Info Cards */}
        <div className={styles.infoCardsContainer}>
          <div className={styles.infoCard}>
            <h3 className={styles.infoCardTitle}>Comandos Disponibles</h3>
            <p className={styles.infoCardContent}>
              mkdisk, rmdisk, fdisk, mount, mounted, mkfs, login, logout, mkgrp, mkusr, mkfile, mkdir, cat, rep
            </p>
          </div>

          <div className={styles.infoCard}>
            <h3 className={styles.infoCardTitle}>Formato de Script</h3>
            <p className={styles.infoCardContent}>
              Los scripts deben tener extensión .smia. Use # para comentarios. Los parámetros se separan por espacios.
            </p>
          </div>

          <div className={styles.infoCard}>
            <h3 className={styles.infoCardTitle}>Estado del Sistema</h3>
            <p className={styles.infoCardContent}>
              {isExecuting ? (
                <span className={styles.statusExecuting}>⚡ Ejecutando comandos...</span>
              ) : (
                <span className={styles.statusReady}>✓ Listo para ejecutar</span>
              )}
            </p>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className={styles.footer}>
        <div className={styles.footerContainer}>
          <p className={styles.footerText}>
            Proyecto 1 - Manejo e Implementación de Archivos | Ingeniería en Ciencias y Sistemas
          </p>
        </div>
      </footer>
    </div>
  );
}

export default App;
