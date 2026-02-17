import Link from 'next/link'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { LoginForm } from '@/components/auth/login-form'
import { RegisterForm } from '@/components/auth/register-form'

export default function LoginPage() {
  return (
    <Card className="border-0 shadow-xl">
      <CardHeader className="text-center pb-2">
        <CardTitle className="text-2xl">Bienvenido a SmartStocks</CardTitle>
        <CardDescription>
          Inicia sesion o crea una cuenta para comenzar
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue="login" className="w-full">
          <TabsList className="grid w-full grid-cols-2 mb-6">
            <TabsTrigger value="login">Iniciar Sesion</TabsTrigger>
            <TabsTrigger value="register">Registrarse</TabsTrigger>
          </TabsList>
          <TabsContent value="login">
            <LoginForm />
            <p className="text-center text-sm text-muted-foreground mt-4">
              <Link href="/forgot-password" className="text-primary hover:underline">
                ¿Olvidaste tu contrasena?
              </Link>
            </p>
          </TabsContent>
          <TabsContent value="register">
            <RegisterForm />
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  )
}
