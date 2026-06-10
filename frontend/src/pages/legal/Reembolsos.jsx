import { Link } from 'react-router-dom'

export default function Reembolsos() {
  return (
    <div className="legal-page">
      <div className="legal-inner">
        <Link to="/" className="back-link">← Volver al inicio</Link>
        <h1 className="legal-title">Política de Reembolso</h1>
        <p className="legal-updated">Última actualización: junio 2026</p>

        <section className="legal-section">
          <h2>Reembolso por cancelación del usuario</h2>
          <p>Podés cancelar tu entrada directamente desde la sección "Mis Entradas" de tu cuenta. Las condiciones son las siguientes:</p>
          <div className="legal-table">
            <div className="legal-table-row legal-table-header">
              <span>Momento de cancelación</span>
              <span>Reembolso</span>
            </div>
            <div className="legal-table-row">
              <span>Más de 48 hs antes del evento</span>
              <span className="legal-ok">100% del valor</span>
            </div>
            <div className="legal-table-row">
              <span>Entre 24 y 48 hs antes</span>
              <span className="legal-warn">50% del valor</span>
            </div>
            <div className="legal-table-row">
              <span>Menos de 24 hs antes</span>
              <span className="legal-no">Sin reembolso</span>
            </div>
          </div>
        </section>

        <section className="legal-section">
          <h2>Reembolso por cancelación del evento</h2>
          <p>Si el organizador cancela el evento, todos los compradores recibirán un reembolso automático del <strong>100% del precio abonado</strong> en un plazo de 5 a 10 días hábiles. Vórtice notificará por email a todos los afectados.</p>
        </section>

        <section className="legal-section">
          <h2>Traspaso de entradas</h2>
          <p>Como alternativa al reembolso, podés traspasar tu entrada a otro usuario registrado en la plataforma de forma gratuita, desde la sección "Mis Entradas". El traspaso es inmediato y no tiene costo adicional.</p>
        </section>

        <section className="legal-section">
          <h2>Cómo se procesa el reembolso</h2>
          <p>Los reembolsos se acreditan en la misma forma de pago utilizada en la compra. El tiempo de acreditación depende de tu banco o entidad emisora, pero generalmente es de 5 a 10 días hábiles desde la aprobación.</p>
        </section>

        <section className="legal-section">
          <h2>Club Vórtice y reembolsos</h2>
          <p>El precio de la membresía mensual del Club Vórtice (<strong>$5.000</strong>) no es reembolsable una vez activada la suscripción, ya que el descuento entra en vigor de inmediato.</p>
        </section>

        <section className="legal-section">
          <h2>Contacto</h2>
          <p>Para dudas o reclamos sobre reembolsos, escribinos a <strong>reembolsos@vortice.com</strong>. Respondemos en un plazo máximo de 48 horas hábiles.</p>
        </section>
      </div>
    </div>
  )
}
