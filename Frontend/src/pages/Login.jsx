import React, { useState } from 'react';
import { HardDrive, User, Lock, LogIn, AlertCircle } from 'lucide-react';
import styles from './Login.module.css';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async () => {
    setError('');

    // Validaciones básicas
    if (!username.trim() || !password.trim()) {
      setError('Por favor, complete todos los campos');
      return;
    }

    setIsLoading(true);

    try {
      // Aquí irá la llamada al backend cuando esté listo
      const response = await fetch('http://localhost:8080/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        const data = await response.json();
        // Guardar token o sesión
        localStorage.setItem('token', data.token);
        // Redirigir al dashboard
        window.location.href = '/dashboard';
      } else {
        setError('Usuario o contraseña incorrectos');
      }
    } catch (err) {
      // Por ahora, simular login exitoso
      console.log('Intento de login:', { username, password });
      setError('Backend no conectado. Usuario: admin, Contraseña: 123');
      
      // Simulación: Si es admin/123, "loguear"
      if (username === 'admin' && password === '123') {
        setTimeout(() => {
          alert('Login exitoso! (Simulación)');
          // Aquí redirigirías al dashboard
        }, 500);
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleSubmit();
    }
  };

  return (
    <div className={styles.loginPage}>
      {/* Efectos de fondo decorativos */}
      <div className={styles.decorativeBg}>
        <div className={styles.bgShape1}></div>
        <div className={styles.bgShape2}></div>
      </div>

      {/* Contenedor principal */}
      <div className={styles.container}>
        {/* Logo y título */}
        <div className={styles.header}>
          <div className={styles.logoContainer}>
            <HardDrive className={styles.logoIcon} />
          </div>
          <h1 className={styles.title}>GoDisk</h1>
          <p className={styles.subtitle}>Sistema de Archivos EXT2</p>
        </div>

        {/* Card de login */}
        <div className={styles.loginCard}>
          <div className={styles.cardHeader}>
            <h2 className={styles.cardTitle}>Iniciar Sesión</h2>
            <p className={styles.cardSubtitle}>Ingrese sus credenciales para continuar</p>
          </div>

          <div className={styles.form}>
            {/* Campo de usuario */}
            <div>
              <label className={styles.inputLabel}>
                Usuario
              </label>
              <div className={styles.inputGroup}>
                <div className={styles.inputIcon}>
                  <User className={styles.icon} />
                </div>
                <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  onKeyPress={handleKeyPress}
                  className={styles.input}
                  placeholder="Ingrese su usuario"
                />
              </div>
            </div>

            {/* Campo de contraseña */}
            <div>
              <label className={styles.inputLabel}>
                Contraseña
              </label>
              <div className={styles.inputGroup}>
                <div className={styles.inputIcon}>
                  <Lock className={styles.icon} />
                </div>
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  onKeyPress={handleKeyPress}
                  className={styles.input}
                  placeholder="Ingrese su contraseña"
                />
              </div>
            </div>

            {/* Mensaje de error */}
            {error && (
              <div className={styles.errorContainer}>
                <AlertCircle className={styles.errorIcon} />
                <p className={styles.errorMessage}>{error}</p>
              </div>
            )}

            {/* Botón de submit */}
            <button
              onClick={handleSubmit}
              disabled={isLoading}
              className={styles.submitButton}
            >
              {isLoading ? (
                <>
                  <div className={styles.spinner}></div>
                  <span>Iniciando sesión...</span>
                </>
              ) : (
                <>
                  <LogIn className={styles.buttonIcon} />
                  <span>Iniciar Sesión</span>
                </>
              )}
            </button>
          </div>

          {/* Enlaces adicionales */}
          <div className={styles.links}>
            <button className={styles.link}>
              ¿Olvidó su contraseña?
            </button>
          </div>
        </div>

        {/* Información adicional */}
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