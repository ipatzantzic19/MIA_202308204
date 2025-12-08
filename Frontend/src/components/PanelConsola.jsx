import { Terminal, Trash2 } from "lucide-react";

function PanelConsola({
  outputCommands,
  handleClearOutput,
  styles,
}) {
  return (
    <div className={styles.card}>
      <div className={styles.cardHeader}>
        <div className={styles.cardHeaderContent}>
          <div className={styles.cardTitleContainer}>
            <Terminal className={`${styles.cardIcon} text-yellow-400`} />
            <h2 className={styles.cardTitle}>Salida de Comandos</h2>
          </div>

          <button onClick={handleClearOutput} className={styles.clearButton}>
            <Trash2 className={styles.buttonIcon} />
            <span>Limpiar</span>
          </button>
        </div>
      </div>

      <div className={styles.cardBody}>
        <div className={styles.outputContainer}>
          <pre className={styles.preformatted}>
            {outputCommands ||
              " Esperando ejecución de comandos...\n Los resultados aparecerán aquí"}
          </pre>
        </div>
      </div>
    </div>
  );
}

export default PanelConsola;