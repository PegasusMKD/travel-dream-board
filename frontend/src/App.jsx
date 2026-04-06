import { Routes, Route } from 'react-router-dom'
import Header from './components/Header'
import Home from './pages/Home'
import BoardDetail from './pages/BoardDetail'

export default function App() {
  return (
    <div className="min-h-screen bg-surface-50 text-text-primary transition-colors duration-200">
      <Header />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-16">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/board/:id" element={<BoardDetail />} />
        </Routes>
      </main>
    </div>
  )
}
