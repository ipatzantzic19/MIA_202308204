import { Terminal, Trash2 } from "lucide-react";
import styles from '../styles/App.module.css';

function PanelConsola({
  outputCommands,       // El string que contiene la salida de los comandos para mostrar.
  handleClearOutput     // La función para limpiar el contenido de la consola.
}) {
  return (
    // Contenedor principal del panel
    <div className={styles.card}>
      {/* Cabecera del panel */}
      <div className={styles.cardHeader}>
        <div className={styles.cardHeaderContent}>
          {/* Título del panel*/}
          <div className={styles.cardTitleContainer}>
            <Terminal className={`${styles.cardIcon} text-yellow-400`} />
            <h2 className={styles.cardTitle}>Salida de Comandos</h2>
          </div>

          {/* Botón para limpiar la salida de la consola */}
          <button onClick={handleClearOutput} className={styles.clearButton}>
            <Trash2 className={styles.buttonIcon} />
            <span>Limpiar</span>
          </button>
        </div>
      </div>

      {/* Cuerpo del panel */}
      <div className={styles.cardBody}>
        {/* Contenedor para el área de salida */}
        <div className={styles.outputContainer}>
          {/* La etiqueta <pre> respeta los espacios en blanco y saltos de línea del texto */}
          <pre className={styles.preformatted}>
            {/* Muestra la salida de los comandos si existe; si no, muestra un mensaje de espera. */}
            {outputCommands || " Esperando ejecución de comandos...\n Los resultados aparecerán aquí"}
          </pre>
        </div>
      </div>
    </div>
  );
}

export default PanelConsola;