import { Link } from 'react-router-dom'

export default function Privacidad() {
  return (
    <div className="legal-page">
      <div className="legal-inner">
        <Link to="/" className="back-link">← Volver al inicio</Link>
        <h1 className="legal-title">Política de Privacidad</h1>
        <p className="legal-updated">Última actualización: junio 2026</p>

        <section className="legal-section">
          <h2>1. Información que recopilamos</h2>
          <p>Recopilamos la información que usted nos proporciona al registrarse: nombre, apellido y dirección de correo electrónico. También registramos las entradas adquiridas y el historial de actividad en la plataforma para brindar un mejor servicio.</p>
        </section>

        <section className="legal-section">
          <h2>2. Uso de la información</h2>
          <p>Utilizamos sus datos para:</p>
          <ul>
            <li>Gestionar su cuenta y las entradas adquiridas.</li>
            <li>Procesar pagos y aplicar descuentos del Club Vórtice.</li>
            <li>Enviar confirmaciones y notificaciones relacionadas con sus compras.</li>
            <li>Mejorar la experiencia en la plataforma.</li>
          </ul>
        </section>

        <section className="legal-section">
          <h2>3. Compartición de datos</h2>
          <p>No vendemos ni compartimos su información personal con terceros, excepto cuando sea necesario para procesar sus compras o cuando lo exija la legislación vigente. Los organizadores de eventos solo reciben los datos mínimos necesarios para validar las entradas.</p>
        </section>

        <section className="legal-section">
          <h2>4. Seguridad</h2>
          <p>Implementamos medidas de seguridad técnicas y organizativas para proteger su información, incluyendo cifrado de contraseñas mediante bcrypt y comunicación segura mediante tokens JWT. Sin embargo, ningún sistema es 100% seguro y no podemos garantizar la seguridad absoluta de los datos.</p>
        </section>

        <section className="legal-section">
          <h2>5. Sus derechos</h2>
          <p>Usted tiene derecho a acceder, rectificar y eliminar sus datos personales. Para ejercer estos derechos, puede contactarnos en <strong>privacidad@vortice.com</strong>. Procesaremos su solicitud en un plazo máximo de 30 días hábiles.</p>
        </section>

        <section className="legal-section">
          <h2>6. Cookies</h2>
          <p>Vórtice no utiliza cookies de rastreo ni publicidad. Solo almacenamos en el navegador la información de sesión necesaria para mantener su inicio de sesión activo (token de autenticación).</p>
        </section>

        <section className="legal-section">
          <h2>7. Contacto</h2>
          <p>Si tiene consultas sobre esta política, puede contactarnos en <strong>privacidad@vortice.com</strong> o mediante el formulario de contacto disponible en la plataforma.</p>
        </section>
      </div>
    </div>
  )
}
