import { Upload, Play, Trash2, Terminal, FileText } from "lucide-react";
import styles from '../styles/App.module.css';

function PanelEditor({
    inputCommands,      // El valor actual del área de texto de comandos.
    setInputCommands,   // Función para actualizar el estado de `inputCommands`.
    fileName,           // El nombre del archivo que se ha cargado.
    fileInputRef,       // Referencia al input de tipo archivo para poder activarlo programáticamente.
    handleFileUpload,   // Controlador para el evento de carga de archivos.
    handleExecute,      // Controlador para el botón de ejecutar comandos.
    handleClearInput,   // Controlador para el botón de limpiar el área de texto.
    isExecuting,        // Estado que indica si se está ejecutando una acción, para deshabilitar botones.
}) {
    return (
        // Contenedor principal del panel del editor
        <div className={styles.card}>
            {/* Cabecera del panel */}
            <div className={styles.cardHeader}>
                <div className={styles.cardHeaderContent}>
                    {/* Contenedor para el título y el ícono. */}
                    <div className={styles.cardTitleContainer}>
                        <Terminal className={`${styles.cardIcon} text-green-400`} />
                        <h2 className={styles.cardTitle}>Entrada de Comandos</h2>
                    </div>

                    {/* Muestra el nombre del archivo si uno ha sido cargado. */}
                    {fileName && (
                        <span className={styles.fileName}>
                            <FileText className={styles.fileNameIcon} />
                            <span>{fileName}</span>
                        </span>
                    )}
                </div>
            </div>

            {/* Cuerpo del panel que contiene el área de texto y los botones. */}
            <div className={styles.cardBody}>
                {/* Área de texto para que el usuario ingrese o vea los comandos. */}
                <textarea
                    value={inputCommands} // El contenido del textarea está controlado por el estado del componente padre.
                    onChange={(e) => setInputCommands(e.target.value)} // Actualiza el estado cuando el usuario escribe.
                    placeholder="# Ingrese sus comandos aquí..."
                    className={styles.textarea}
                    spellCheck="false" // Deshabilita la corrección ortográfica, útil para código.
                />

                {/* Contenedor para los botones de acción. */}
                <div className={styles.buttonGroup}>
                    {/* Input de archivo oculto. Se activa a través del botón "Cargar Script". */}
                    <input
                        type="file"
                        ref={fileInputRef} // Se asigna la referencia para poder acceder a este elemento.
                        onChange={handleFileUpload} // Función que se llama al seleccionar un archivo.
                        accept=".mia" // Filtra los archivos para mostrar solo la extensión .smia.
                        className={styles.hiddenInput}
                    />

                    {/* Botón para activar el input de archivo y cargar un script. */}
                    <button
                        onClick={() => fileInputRef.current?.click()} // Al hacer clic, se simula un clic en el input de archivo.
                        className={`${styles.button} ${styles.buttonPrimary}`}
                    >
                        <Upload className={styles.buttonIcon} />
                        <span>Cargar Script</span>
                    </button>

                    {/* Botón para ejecutar los comandos del área de texto. */}
                    <button
                        onClick={handleExecute}
                        // Se deshabilita si hay una ejecución en curso o si el área de texto está vacía.
                        disabled={isExecuting || !inputCommands.trim()}
                        className={`${styles.button} ${styles.buttonSuccess}`}
                    >
                        <Play className={styles.buttonIcon} />
                        {/* El texto del botón cambia para indicar que se está ejecutando. */}
                        <span>{isExecuting ? "Ejecutando..." : "Ejecutar"}</span>
                    </button>

                    {/* Botón para borrar el contenido del área de texto. */}
                    <button
                        onClick={handleClearInput}
                        className={styles.buttonSecondary}
                        title="Limpiar entrada" // Texto que aparece al pasar el cursor sobre el botón.
                    >
                        <Trash2 className={styles.buttonIcon} />
                    </button>
                </div>
            </div>
        </div>
    );
}

export default PanelEditor;