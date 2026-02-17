'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import {
  User,
  Mail,
  School,
  Camera,
  Bell,
  Shield,
  LogOut,
  Loader2,
  Save,
  Check,
} from 'lucide-react'
import { toast } from 'sonner'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Separator } from '@/components/ui/separator'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { useAuthStore } from '@/lib/stores/auth-store'
import { userApi, authApi } from '@/lib/api/auth'
import type { School as SchoolType } from '@/types'

export default function SettingsPage() {
  const router = useRouter()
  const { user, updateUser, logout, refreshToken } = useAuthStore()

  const [schools, setSchools] = useState<SchoolType[]>([])
  const [loading, setLoading] = useState(false)
  const [saving, setSaving] = useState(false)

  // Form state
  const [username, setUsername] = useState('')
  const [profilePictureUrl, setProfilePictureUrl] = useState('')
  const [schoolId, setSchoolId] = useState('')

  // Notification preferences (stored in localStorage)
  const [notifications, setNotifications] = useState({
    pvpInvites: true,
    tournamentReminders: true,
    rankChanges: true,
    dailyReminder: false,
  })

  useEffect(() => {
    if (user) {
      setUsername(user.username)
      setProfilePictureUrl(user.profile_picture_url || '')
      setSchoolId(user.school_id || '')
    }

    // Load notification preferences from localStorage
    const savedNotifications = localStorage.getItem('notification_preferences')
    if (savedNotifications) {
      setNotifications(JSON.parse(savedNotifications))
    }

    // Fetch schools
    fetchSchools()
  }, [user])

  const fetchSchools = async () => {
    try {
      const data = await authApi.getSchools()
      setSchools(data)
    } catch (error) {
      console.error('Error fetching schools:', error)
    }
  }

  const handleSaveProfile = async () => {
    if (!username.trim()) {
      toast.error('El nombre de usuario no puede estar vacio')
      return
    }

    setSaving(true)
    try {
      const updatedUser = await userApi.updateProfile({
        username: username.trim(),
        profile_picture_url: profilePictureUrl.trim() || undefined,
        school_id: schoolId || undefined,
      })

      updateUser(updatedUser)
      toast.success('Perfil actualizado correctamente')
    } catch (error: any) {
      toast.error(error.response?.data?.message || 'Error al actualizar el perfil')
    } finally {
      setSaving(false)
    }
  }

  const handleNotificationChange = (key: keyof typeof notifications, value: boolean) => {
    const newNotifications = { ...notifications, [key]: value }
    setNotifications(newNotifications)
    localStorage.setItem('notification_preferences', JSON.stringify(newNotifications))
    toast.success('Preferencias guardadas')
  }

  const handleLogout = async () => {
    try {
      if (refreshToken) {
        await authApi.logout(refreshToken)
      }
    } catch (error) {
      // Continue with logout even if API call fails
    }
    logout()
    router.push('/login')
  }

  if (!user) {
    return (
      <div className="flex items-center justify-center min-h-[60vh]">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    )
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <div>
        <h1 className="text-2xl font-bold">Configuracion</h1>
        <p className="text-muted-foreground">Administra tu cuenta y preferencias</p>
      </div>

      {/* Profile Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <User className="w-5 h-5" />
            Perfil
          </CardTitle>
          <CardDescription>
            Actualiza tu informacion personal
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          {/* Avatar */}
          <div className="flex items-center gap-4">
            <Avatar className="w-20 h-20">
              <AvatarImage src={profilePictureUrl || undefined} />
              <AvatarFallback className="bg-primary text-white text-2xl">
                {username.charAt(0).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div className="flex-1">
              <Label htmlFor="profilePicture">URL de foto de perfil</Label>
              <div className="flex gap-2 mt-1">
                <Input
                  id="profilePicture"
                  placeholder="https://ejemplo.com/mi-foto.jpg"
                  value={profilePictureUrl}
                  onChange={(e) => setProfilePictureUrl(e.target.value)}
                />
              </div>
              <p className="text-xs text-muted-foreground mt-1">
                Ingresa la URL de una imagen para tu perfil
              </p>
            </div>
          </div>

          <Separator />

          {/* Username */}
          <div className="space-y-2">
            <Label htmlFor="username">Nombre de usuario</Label>
            <Input
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Tu nombre de usuario"
            />
          </div>

          {/* Email (read-only) */}
          <div className="space-y-2">
            <Label htmlFor="email">Correo electronico</Label>
            <div className="flex items-center gap-2">
              <Input
                id="email"
                value={user.email}
                disabled
                className="bg-gray-50"
              />
              {user.email_verified ? (
                <span className="flex items-center text-sm text-green-600">
                  <Check className="w-4 h-4 mr-1" />
                  Verificado
                </span>
              ) : (
                <span className="text-sm text-yellow-600">No verificado</span>
              )}
            </div>
          </div>

          {/* School */}
          <div className="space-y-2">
            <Label htmlFor="school">Colegio</Label>
            <Select value={schoolId || "none"} onValueChange={(val) => setSchoolId(val === "none" ? "" : val)}>
              <SelectTrigger>
                <SelectValue placeholder="Selecciona tu colegio" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">Sin colegio</SelectItem>
                {schools.map((school) => (
                  <SelectItem key={school.id} value={school.id}>
                    {school.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <p className="text-xs text-muted-foreground">
              Asocia tu cuenta a un colegio para competir en rankings escolares
            </p>
          </div>

          <Button onClick={handleSaveProfile} disabled={saving} className="w-full">
            {saving ? (
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
            ) : (
              <Save className="w-4 h-4 mr-2" />
            )}
            Guardar cambios
          </Button>
        </CardContent>
      </Card>

      {/* Notification Settings */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Bell className="w-5 h-5" />
            Notificaciones
          </CardTitle>
          <CardDescription>
            Configura tus preferencias de notificaciones
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">Invitaciones PvP</p>
              <p className="text-sm text-muted-foreground">
                Recibir notificaciones de invitaciones a batallas
              </p>
            </div>
            <Switch
              checked={notifications.pvpInvites}
              onCheckedChange={(value) => handleNotificationChange('pvpInvites', value)}
            />
          </div>

          <Separator />

          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">Recordatorios de torneos</p>
              <p className="text-sm text-muted-foreground">
                Notificaciones sobre torneos proximos
              </p>
            </div>
            <Switch
              checked={notifications.tournamentReminders}
              onCheckedChange={(value) => handleNotificationChange('tournamentReminders', value)}
            />
          </div>

          <Separator />

          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">Cambios de rango</p>
              <p className="text-sm text-muted-foreground">
                Notificaciones cuando subas o bajes de rango
              </p>
            </div>
            <Switch
              checked={notifications.rankChanges}
              onCheckedChange={(value) => handleNotificationChange('rankChanges', value)}
            />
          </div>

          <Separator />

          <div className="flex items-center justify-between">
            <div>
              <p className="font-medium">Recordatorio diario</p>
              <p className="text-sm text-muted-foreground">
                Recordatorio para completar tu quiz diario
              </p>
            </div>
            <Switch
              checked={notifications.dailyReminder}
              onCheckedChange={(value) => handleNotificationChange('dailyReminder', value)}
            />
          </div>
        </CardContent>
      </Card>

      {/* Account Actions */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="w-5 h-5" />
            Cuenta
          </CardTitle>
          <CardDescription>
            Acciones relacionadas con tu cuenta
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="destructive" className="w-full">
                <LogOut className="w-4 h-4 mr-2" />
                Cerrar sesion
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Cerrar sesion</AlertDialogTitle>
                <AlertDialogDescription>
                  Estas seguro de que quieres cerrar sesion? Tendras que volver a iniciar sesion para acceder a tu cuenta.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancelar</AlertDialogCancel>
                <AlertDialogAction onClick={handleLogout}>
                  Cerrar sesion
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </CardContent>
      </Card>
    </div>
  )
}
