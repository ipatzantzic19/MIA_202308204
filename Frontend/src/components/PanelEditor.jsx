import { Upload, Play, Trash2, Terminal, FileText } from "lucide-react";
import styles from '../styles/App.module.css';

function PanelEditor({
    inputCommands,
    setInputCommands,
    fileName,
    fileInputRef,
    handleFileUpload,
    handleExecute,
    handleClearInput,
    isExecuting,
}) {
    return (
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
                    placeholder=" Ingrese sus comandos aquí..."
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
                        <span>{isExecuting ? "Ejecutando..." : "Ejecutar"}</span>
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
    );
}

export default PanelEditor;