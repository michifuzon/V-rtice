import { Link } from 'react-router-dom'

export default function Footer() {
  const year = new Date().getFullYear()

  return (
    <footer className="footer">
      <div className="footer-inner">
        <div>
          <p className="footer-brand-name">Vórtice</p>
          <p className="footer-brand-desc">Eventos que atrapan.</p>
        </div>

        <div className="footer-cols">
          <div>
            <p className="footer-col-title">Navegación</p>
            <ul className="footer-links">
              <li><Link to="/">Eventos</Link></li>
              <li><Link to="/mis-entradas">Mis entradas</Link></li>
              <li><Link to="/mis-favoritos">Mis favoritos</Link></li>
              <li><Link to="/register">Crear cuenta</Link></li>
              <li><Link to="/login">Iniciar sesión</Link></li>
            </ul>
          </div>

          <div>
            <p className="footer-col-title">Legal</p>
            <ul className="footer-links">
              <li><Link to="/terminos">Términos y condiciones</Link></li>
              <li><Link to="/privacidad">Política de privacidad</Link></li>
              <li><Link to="/reembolsos">Política de reembolso</Link></li>
            </ul>
          </div>
        </div>
      </div>

      <div className="footer-bottom">
        <span>© {year} Vórtice. Todos los derechos reservados.</span>
        <span>Trabajo final — Desarrollo de Software {year}</span>
      </div>
    </footer>
  )
}
