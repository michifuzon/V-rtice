import { Link } from 'react-router-dom'

export default function Terminos() {
  return (
    <div className="legal-page">
      <div className="legal-inner">
        <Link to="/" className="back-link">← Volver al inicio</Link>
        <h1 className="legal-title">Términos y Condiciones</h1>
        <p className="legal-updated">Última actualización: junio 2026</p>

        <section className="legal-section">
          <h2>1. Aceptación de los términos</h2>
          <p>Al acceder y utilizar la plataforma Vórtice, el usuario acepta quedar vinculado por los presentes Términos y Condiciones. Si no está de acuerdo con alguno de estos términos, le pedimos que no utilice el servicio.</p>
        </section>

        <section className="legal-section">
          <h2>2. Descripción del servicio</h2>
          <p>Vórtice es una plataforma de venta de entradas para eventos culturales, deportivos y de entretenimiento. Actuamos como intermediarios entre los organizadores de eventos y el público, facilitando la adquisición de entradas de forma segura y conveniente.</p>
        </section>

        <section className="legal-section">
          <h2>3. Registro y cuenta de usuario</h2>
          <p>Para adquirir entradas es necesario crear una cuenta. El usuario es responsable de mantener la confidencialidad de sus credenciales y de todas las actividades realizadas bajo su cuenta. Vórtice no será responsable de pérdidas causadas por el uso no autorizado de su cuenta.</p>
        </section>

        <section className="legal-section">
          <h2>4. Compra de entradas</h2>
          <p>Todas las compras son definitivas una vez confirmadas. El precio de las entradas incluye los impuestos aplicables y el cargo por servicio de la plataforma. Vórtice se reserva el derecho de modificar los precios sin previo aviso.</p>
        </section>

        <section className="legal-section">
          <h2>5. Club Vórtice</h2>
          <p>La membresía del Club Vórtice es una suscripción mensual de $5.000 que otorga un 10% de descuento en todas las entradas adquiridas mientras la suscripción esté activa. La suscripción se renueva manualmente y no se realiza ningún cobro automático. El descuento no es acumulable con otras promociones.</p>
        </section>

        <section className="legal-section">
          <h2>6. Cancelación y reembolsos</h2>
          <p>Las entradas pueden cancelarse desde la plataforma hasta 48 horas antes del evento. Los reembolsos se procesan en un plazo de 5 a 10 días hábiles. En caso de cancelación del evento por parte del organizador, se reembolsará el 100% del valor abonado.</p>
        </section>

        <section className="legal-section">
          <h2>7. Limitación de responsabilidad</h2>
          <p>Vórtice no se responsabiliza por la calidad, seguridad o desarrollo de los eventos. La responsabilidad de Vórtice se limita al valor de las entradas adquiridas a través de la plataforma.</p>
        </section>

        <section className="legal-section">
          <h2>8. Modificaciones</h2>
          <p>Vórtice se reserva el derecho de modificar estos términos en cualquier momento. Los cambios entrarán en vigor desde su publicación en la plataforma. El uso continuado del servicio implica la aceptación de los nuevos términos.</p>
        </section>
      </div>
    </div>
  )
}
