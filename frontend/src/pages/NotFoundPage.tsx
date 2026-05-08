import { Link } from 'react-router-dom'

export const NotFoundPage = () => {
  return (
    <div className="page-body">
      <section className="card">
        <span className="kicker">Lost in the schedule</span>
        <h2>That page does not exist.</h2>
        <p>Head back to the dashboard to keep working.</p>
        <Link to="/dashboard" className="button inline">
          Go to dashboard
        </Link>
      </section>
    </div>
  )
}
