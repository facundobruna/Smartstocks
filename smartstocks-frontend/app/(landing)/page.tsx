import Link from 'next/link'
import { ArrowRight, BookOpen, TrendingUp, Target, MessageCircle, Users, GraduationCap, Award } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'

const stats = [
  { value: '10,000+', label: 'Estudiantes activos' },
  { value: '500+', label: 'Lecciones disponibles' },
  { value: '95%', label: 'Satisfaccion estudiantil' },
]

const features = [
  {
    icon: BookOpen,
    title: 'Educacion',
    description: 'Aprende los conceptos fundamentales de finanzas personales e inversiones de forma interactiva.',
    href: '/education',
  },
  {
    icon: TrendingUp,
    title: 'Simulador',
    description: 'Practica tus habilidades de trading con escenarios reales sin arriesgar dinero real.',
    href: '/simulator',
  },
  {
    icon: Target,
    title: 'Nosotros',
    description: 'Conoce nuestra mision de democratizar la educacion financiera para jovenes argentinos.',
    href: '/about',
  },
  {
    icon: MessageCircle,
    title: 'Contacto',
    description: '¿Tienes preguntas? Nuestro equipo esta listo para ayudarte en tu camino financiero.',
    href: '/contact',
  },
]

const benefits = [
  {
    icon: Users,
    title: 'Comunidad Activa',
    description: 'Conecta con otros jovenes interesados en finanzas y comparte conocimientos.',
  },
  {
    icon: GraduationCap,
    title: 'Contenido Actualizado',
    description: 'Aprende con material desarrollado por expertos y actualizado constantemente.',
  },
  {
    icon: Award,
    title: 'Gamificacion',
    description: 'Gana puntos, sube de rango y desbloquea logros mientras aprendes.',
  },
]

export default function LandingPage() {
  return (
    <>
      {/* Hero Section */}
      <section className="relative bg-gradient-to-b from-blue-50 to-white py-20 lg:py-32">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto text-center">
            <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold text-foreground mb-6">
              SMARTSTOCKS
            </h1>
            <p className="text-xl md:text-2xl text-muted-foreground mb-8">
              Aprender hoy lo que cambiara tu manana
            </p>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto mb-10">
              La plataforma de educacion financiera gamificada disenada para jovenes argentinos.
              Aprende a invertir de forma divertida y segura.
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center mb-16">
              <Button size="lg" asChild className="text-lg px-8">
                <Link href="/register">
                  Comenzar Gratis
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button size="lg" variant="outline" asChild className="text-lg px-8">
                <Link href="/about">Ver Demo</Link>
              </Button>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {stats.map((stat) => (
                <Card key={stat.label} className="bg-white/80 backdrop-blur-sm border-0 shadow-lg">
                  <CardContent className="pt-6 text-center">
                    <p className="text-4xl font-bold text-primary mb-2">{stat.value}</p>
                    <p className="text-sm text-muted-foreground">{stat.label}</p>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-white">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
              Todo lo que necesitas para aprender finanzas
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Herramientas disenadas para que aprender sobre dinero e inversiones sea facil y divertido.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {features.map((feature) => {
              const Icon = feature.icon
              return (
                <Link key={feature.title} href={feature.href}>
                  <Card className="h-full transition-all duration-300 hover:-translate-y-1 hover:shadow-lg cursor-pointer group">
                    <CardContent className="pt-8 pb-6 text-center">
                      <div className="w-16 h-16 mx-auto mb-6 rounded-2xl bg-primary/10 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
                        <Icon className="w-8 h-8 text-primary" />
                      </div>
                      <h3 className="text-xl font-semibold text-foreground mb-3">
                        {feature.title}
                      </h3>
                      <p className="text-sm text-muted-foreground">
                        {feature.description}
                      </p>
                    </CardContent>
                  </Card>
                </Link>
              )
            })}
          </div>
        </div>
      </section>

      {/* Benefits Section */}
      <section className="py-20 bg-gray-50">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
              ¿Por que SmartStocks?
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Creemos que todos los jovenes merecen acceso a educacion financiera de calidad.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-5xl mx-auto">
            {benefits.map((benefit) => {
              const Icon = benefit.icon
              return (
                <div key={benefit.title} className="text-center">
                  <div className="w-14 h-14 mx-auto mb-4 rounded-full bg-primary/10 flex items-center justify-center">
                    <Icon className="w-7 h-7 text-primary" />
                  </div>
                  <h3 className="text-lg font-semibold text-foreground mb-2">
                    {benefit.title}
                  </h3>
                  <p className="text-sm text-muted-foreground">
                    {benefit.description}
                  </p>
                </div>
              )
            })}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-primary">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
            Comienza tu viaje financiero hoy
          </h2>
          <p className="text-lg text-white/80 max-w-2xl mx-auto mb-8">
            Unete a miles de jovenes que ya estan aprendiendo a manejar su dinero de forma inteligente.
          </p>
          <Button size="lg" variant="secondary" asChild className="text-lg px-8">
            <Link href="/register">
              Crear Cuenta Gratis
              <ArrowRight className="ml-2 h-5 w-5" />
            </Link>
          </Button>
        </div>
      </section>
    </>
  )
}
