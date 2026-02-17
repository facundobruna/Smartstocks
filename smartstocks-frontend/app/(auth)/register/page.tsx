import Link from 'next/link'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { LoginForm } from '@/components/auth/login-form'
import { RegisterForm } from '@/components/auth/register-form'

export default function RegisterPage() {
  return (
    <Card className="border-0 shadow-xl">
      <CardHeader className="text-center pb-2">
        <CardTitle className="text-2xl">Unete a SmartStocks</CardTitle>
        <CardDescription>
          Crea tu cuenta y comienza a aprender finanzas
        </CardDescription>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue="register" className="w-full">
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
