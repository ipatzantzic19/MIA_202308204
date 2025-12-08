import styles from '../styles/App.module.css';

function InfoCards({ Card }) {
    return (
        <div className={styles.infoCard}>
            <h3 className={styles.infoCardTitle}>{Card.title}</h3>
            <p className={styles.infoCardContent}>{Card.content}</p>
        </div>
    );
}

export default InfoCards;