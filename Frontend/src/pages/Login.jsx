import React, { useState } from 'react';
import { HardDrive, User, Lock, LogIn, AlertCircle } from 'lucide-react';
import styles from '../styles/Login.module.css';

export default function LoginPage() {

  // Estado para almacenar el nombre de usuario ingresado.
  const [username, setUsername] = useState('');
  // Estado para almacenar la contraseña ingresada.
  const [password, setPassword] = useState('');
  // Estado para almacenar mensajes de error y mostrarlos al usuario.
  const [error, setError] = useState('');
  // Estado para gestionar la visualización de indicadores de carga durante la comunicación con el backend.
  const [isLoading, setIsLoading] = useState(false);

  /**
   * Función asincrónica `handleSubmit` que se ejecuta al intentar iniciar sesión.
   * Realiza validaciones, envía las credenciales al backend y maneja la respuesta.
   */
  const handleSubmit = async () => {
    // Limpia cualquier error previo al iniciar un nuevo intento.
    setError('');

    // Validación para asegurar que ambos campos estén completos.
    if (!username.trim() || !password.trim()) {
      setError('Por favor, complete todos los campos');
      return; // Detiene la ejecución si la validación falla.
    }

    // Activa el estado de carga para dar feedback visual al usuario.
    setIsLoading(true);

    try {
      // Intento de conexión con el endpoint del backend para la autenticación.
      const response = await fetch('http://localhost:8080/api/login', {
        method: 'POST', // Método HTTP para enviar datos.
        headers: {
          'Content-Type': 'application/json', // Especifica que el cuerpo es JSON.
        },
        // Convierte los datos del formulario a una cadena JSON.
        body: JSON.stringify({ username, password }),
      });

      // Si la respuesta del servidor es exitosa (ej. status 200-299).
      if (response.ok) {
        const data = await response.json(); // Parsea la respuesta JSON.
        // Almacena el token de sesión en el almacenamiento local del navegador.
        localStorage.setItem('token', data.token);
        // Redirige al usuario al dashboard principal de la aplicación.
        window.location.href = '/dashboard';
      } else {
        // Si las credenciales son incorrectas o hay otro error del servidor.
        setError('Usuario o contraseña incorrectos');
      }
    } catch (err) {
      // --- BLOQUE DE SIMULACIÓN (para desarrollo sin backend) ---
      // Este bloque se ejecuta si la llamada `fetch` falla (ej. backend no disponible).
      console.log('Intento de login (simulación):', { username, password });
      setError('Backend no conectado. Use: admin / 123 para simular.');
      
      // Lógica de simulación para permitir el acceso con credenciales predefinidas.
      if (username === 'admin' && password === '123') {
        setTimeout(() => {
          alert('Login exitoso! (Simulación)');
          // En un caso real, aquí también se redirigiría al dashboard.
        }, 500);
      }
    } finally {
      // Desactiva el estado de carga, independientemente del resultado (éxito o error).
      setIsLoading(false);
    }
  };

  /**
   * Función `handleKeyPress` que permite iniciar sesión al presionar la tecla "Enter".
   * @param {React.KeyboardEvent} e - El evento de teclado.
   */
  const handleKeyPress = (e) => {
    // Si la tecla presionada es "Enter", llama a la función de submit.
    if (e.key === 'Enter') {
      handleSubmit();
    }
  };

  // --- RENDERIZADO DEL COMPONENTE ---
  return (
    // Contenedor principal de la página con su clase de estilo.
    <div className={styles.loginPage}>
      {/* Elementos decorativos para el fondo de la página. */}
      <div className={styles.decorativeBg}>
        <div className={styles.bgShape1}></div>
        <div className={styles.bgShape2}></div>
      </div>

      {/* Contenedor centrado que agrupa todo el contenido principal. */}
      <div className={styles.container}>
        {/* Cabecera con el logo y título de la aplicación. */}
        <div className={styles.header}>
          <div className={styles.logoContainer}>
            <HardDrive className={styles.logoIcon} />
          </div>
          <h1 className={styles.title}>GoDisk</h1>
          <p className={styles.subtitle}>Sistema de Archivos EXT2</p>
        </div>

        {/* Tarjeta (card) que contiene el formulario de inicio de sesión. */}
        <div className={styles.loginCard}>
          <div className={styles.cardHeader}>
            <h2 className={styles.cardTitle}>Iniciar Sesión</h2>
            <p className={styles.cardSubtitle}>Ingrese sus credenciales para continuar</p>
          </div>

          {/* Formulario de inicio de sesión. */}
          <div className={styles.form}>
            {/* Campo de entrada para el nombre de usuario. */}
            <div>
              <label className={styles.inputLabel}>Usuario</label>
              <div className={styles.inputGroup}>
                <div className={styles.inputIcon}>
                  <User className={styles.icon} />
                </div>
                <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)} // Actualiza el estado en cada cambio.
                  onKeyPress={handleKeyPress} // Asocia la función para la tecla "Enter".
                  className={styles.input}
                  placeholder="Ingrese su usuario"
                />
              </div>
            </div>

            {/* Campo de entrada para la contraseña. */}
            <div>
              <label className={styles.inputLabel}>Contraseña</label>
              <div className={styles.inputGroup}>
                <div className={styles.inputIcon}>
                  <Lock className={styles.icon} />
                </div>
                <input
                  type="password" 
                  value={password}
                  onChange={(e) => setPassword(e.target.value)} // Actualiza el estado en cada cambio.
                  onKeyPress={handleKeyPress} // Asocia la función para la tecla "Enter".
                  className={styles.input}
                  placeholder="Ingrese su contraseña"
                />
              </div>
            </div>

            {/* Contenedor para mostrar el mensaje de error (solo si `error` no está vacío). */}
            {error && (
              <div className={styles.errorContainer}>
                <AlertCircle className={styles.errorIcon} />
                <p className={styles.errorMessage}>{error}</p>
              </div>
            )}

            {/* Botón para enviar el formulario. */}
            <button
              onClick={handleSubmit}
              disabled={isLoading} // El botón se deshabilita durante la carga.
              className={styles.submitButton}
            >
              {isLoading ? (
                // Muestra un spinner y texto de carga si `isLoading` es true.
                <>
                  <div className={styles.spinner}></div>
                  <span>Iniciando sesión...</span>
                </>
              ) : (
                // Muestra el icono y texto por defecto.
                <>
                  <LogIn className={styles.buttonIcon} />
                  <span>Iniciar Sesión</span>
                </>
              )}
            </button>
          </div>
        </div>

        {/* Pie de página */}
        <div className={styles.footer}>
          <p className={styles.footerText}>
            Universidad San Carlos de Guatemala
          </p>
          <p className={styles.footerSubtext}>
            MIA - Proyecto 1 | Ingeniería en Ciencias y Sistemas
          </p>
        </div>
      </div>
    </div>
  );
}