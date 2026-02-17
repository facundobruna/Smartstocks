'use client'

import Link from 'next/link'
import {
  TrendingUp,
  Users,
  Target,
  Award,
  BookOpen,
  Gamepad2,
  Shield,
  Heart,
  ArrowRight,
  CheckCircle,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Navbar } from '@/components/layout/navbar'
import { Footer } from '@/components/layout/footer'

const teamMembers = [
  {
    name: 'Equipo SmartStocks',
    role: 'Desarrollo & Educación',
    description: 'Un equipo apasionado por democratizar la educación financiera para los jóvenes argentinos.',
  },
]

const values = [
  {
    icon: BookOpen,
    title: 'Educación Accesible',
    description: 'Creemos que todos los jóvenes merecen acceso a educación financiera de calidad, sin importar su situación económica.',
  },
  {
    icon: Gamepad2,
    title: 'Aprendizaje Divertido',
    description: 'Transformamos conceptos financieros complejos en experiencias de juego atractivas y memorables.',
  },
  {
    icon: Shield,
    title: 'Ambiente Seguro',
    description: 'Proporcionamos un entorno libre de riesgos donde los errores son oportunidades de aprendizaje.',
  },
  {
    icon: Heart,
    title: 'Impacto Social',
    description: 'Buscamos generar un cambio positivo en la sociedad, preparando a la próxima generación para un futuro financiero sólido.',
  },
]

const milestones = [
  { number: '10,000+', label: 'Estudiantes activos' },
  { number: '500+', label: 'Escuelas participantes' },
  { number: '1M+', label: 'Simulaciones completadas' },
  { number: '95%', label: 'Satisfacción de usuarios' },
]

const features = [
  'Simulador de inversiones sin riesgo real',
  'Modo PvP para competir con amigos',
  'Torneos semanales y mensuales',
  'Rankings por escuela y nacional',
  'Contenido educativo gamificado',
  'Sistema de logros y recompensas',
]

export default function AboutPage() {
  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />

      <main className="flex-1">
        {/* Hero Section */}
        <section className="bg-gradient-to-br from-blue-600 via-blue-700 to-purple-700 text-white py-20">
          <div className="container mx-auto px-4 text-center">
            <div className="inline-flex items-center justify-center w-20 h-20 rounded-full bg-white/10 mb-6">
              <TrendingUp className="w-10 h-10" />
            </div>
            <h1 className="text-4xl md:text-5xl font-bold mb-6">
              Sobre SmartStocks
            </h1>
            <p className="text-xl text-blue-100 max-w-3xl mx-auto">
              Somos una plataforma educativa que enseña a los jóvenes argentinos
              sobre finanzas e inversiones a través de simulaciones gamificadas y
              competencias emocionantes.
            </p>
          </div>
        </section>

        {/* Mission Section */}
        <section className="py-16 bg-white">
          <div className="container mx-auto px-4">
            <div className="max-w-4xl mx-auto">
              <div className="grid md:grid-cols-2 gap-12 items-center">
                <div>
                  <h2 className="text-3xl font-bold mb-4">Nuestra Misión</h2>
                  <p className="text-gray-600 mb-4">
                    Democratizar la educación financiera para los jóvenes de Argentina,
                    brindándoles las herramientas y conocimientos necesarios para tomar
                    decisiones financieras inteligentes desde temprana edad.
                  </p>
                  <p className="text-gray-600">
                    Creemos que la alfabetización financiera es una habilidad esencial
                    que debería ser accesible para todos, no solo para quienes tienen
                    recursos económicos.
                  </p>
                </div>
                <div>
                  <h2 className="text-3xl font-bold mb-4">Nuestra Visión</h2>
                  <p className="text-gray-600 mb-4">
                    Ser la plataforma líder de educación financiera gamificada en
                    Latinoamérica, preparando a millones de jóvenes para un futuro
                    financiero exitoso.
                  </p>
                  <p className="text-gray-600">
                    Imaginamos un mundo donde cada joven tiene la confianza y el
                    conocimiento para manejar sus finanzas personales e inversiones
                    de manera responsable.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Stats Section */}
        <section className="py-16 bg-gray-50">
          <div className="container mx-auto px-4">
            <h2 className="text-3xl font-bold text-center mb-12">Nuestro Impacto</h2>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-8 max-w-4xl mx-auto">
              {milestones.map((milestone, index) => (
                <div key={index} className="text-center">
                  <div className="text-4xl md:text-5xl font-bold text-primary mb-2">
                    {milestone.number}
                  </div>
                  <div className="text-gray-600">{milestone.label}</div>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Values Section */}
        <section className="py-16 bg-white">
          <div className="container mx-auto px-4">
            <h2 className="text-3xl font-bold text-center mb-12">Nuestros Valores</h2>
            <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-6xl mx-auto">
              {values.map((value, index) => (
                <Card key={index} className="text-center">
                  <CardContent className="pt-6">
                    <div className="w-14 h-14 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-4">
                      <value.icon className="w-7 h-7 text-primary" />
                    </div>
                    <h3 className="font-semibold text-lg mb-2">{value.title}</h3>
                    <p className="text-sm text-gray-600">{value.description}</p>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section className="py-16 bg-gray-50">
          <div className="container mx-auto px-4">
            <div className="max-w-4xl mx-auto">
              <h2 className="text-3xl font-bold text-center mb-12">
                ¿Qué ofrecemos?
              </h2>
              <div className="grid md:grid-cols-2 gap-4">
                {features.map((feature, index) => (
                  <div key={index} className="flex items-center gap-3 p-4 bg-white rounded-lg">
                    <CheckCircle className="w-5 h-5 text-green-500 flex-shrink-0" />
                    <span>{feature}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </section>

        {/* For Schools Section */}
        <section className="py-16 bg-white">
          <div className="container mx-auto px-4">
            <div className="max-w-4xl mx-auto text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-orange-100 mb-6">
                <Award className="w-8 h-8 text-orange-500" />
              </div>
              <h2 className="text-3xl font-bold mb-4">Para Escuelas</h2>
              <p className="text-gray-600 mb-8 max-w-2xl mx-auto">
                Ofrecemos programas especiales para instituciones educativas.
                Integrá SmartStocks en tu currículo y dale a tus estudiantes
                una ventaja competitiva en educación financiera.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Button size="lg" asChild>
                  <Link href="/register">
                    Registrar mi Escuela
                    <ArrowRight className="ml-2 w-4 h-4" />
                  </Link>
                </Button>
                <Button size="lg" variant="outline" asChild>
                  <a href="mailto:escuelas@smartstocks.com.ar">
                    Contactar Ventas
                  </a>
                </Button>
              </div>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="py-16 bg-gradient-to-r from-primary to-blue-600 text-white">
          <div className="container mx-auto px-4 text-center">
            <h2 className="text-3xl font-bold mb-4">
              ¿Listo para empezar tu viaje financiero?
            </h2>
            <p className="text-blue-100 mb-8 max-w-xl mx-auto">
              Únete a miles de estudiantes que ya están aprendiendo a invertir
              de manera inteligente con SmartStocks.
            </p>
            <Button size="lg" variant="secondary" asChild>
              <Link href="/register">
                Crear Cuenta Gratis
                <ArrowRight className="ml-2 w-4 h-4" />
              </Link>
            </Button>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  )
}
