'use client'

import { useState, useEffect } from 'react'
import {
  MessageSquare,
  ThumbsUp,
  MessageCircle,
  Eye,
  Clock,
  TrendingUp,
  HelpCircle,
  Lightbulb,
  Search,
  Plus,
  Pin,
  Flame,
  Users,
} from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/lib/stores/auth-store'
import { toast } from 'sonner'

interface ForumPost {
  id: string
  title: string
  content: string
  author: {
    username: string
    avatar?: string
    rank: string
  }
  category: 'general' | 'strategies' | 'help' | 'news'
  likes: number
  replies: number
  views: number
  isPinned?: boolean
  isHot?: boolean
  createdAt: string
}

const mockPosts: ForumPost[] = [
  {
    id: '1',
    title: 'Guía para principiantes: Cómo empezar en el simulador',
    content: 'Hola a todos! Creé esta guía para ayudar a los nuevos usuarios a entender...',
    author: { username: 'TraderPro', rank: 'Oro 1' },
    category: 'help',
    likes: 156,
    replies: 42,
    views: 1250,
    isPinned: true,
    createdAt: '2024-01-08',
  },
  {
    id: '2',
    title: 'Mi estrategia para ganar en modo PvP - 80% win rate',
    content: 'Después de 100 partidas, quiero compartir lo que aprendí sobre el modo PvP...',
    author: { username: 'InversorNinja', rank: 'Maestro' },
    category: 'strategies',
    likes: 234,
    replies: 89,
    views: 3420,
    isHot: true,
    createdAt: '2024-01-09',
  },
  {
    id: '3',
    title: 'Nuevo torneo mensual anunciado - Premio de 10,000 tokens!',
    content: 'Acaban de anunciar el torneo de enero con premios increíbles...',
    author: { username: 'NewsBot', rank: 'Plata 2' },
    category: 'news',
    likes: 89,
    replies: 23,
    views: 890,
    createdAt: '2024-01-10',
  },
  {
    id: '4',
    title: '¿Cómo interpretar las noticias del simulador?',
    content: 'Tengo dudas sobre cómo analizar las noticias antes de tomar una decisión...',
    author: { username: 'NuevoTrader', rank: 'Bronce 2' },
    category: 'help',
    likes: 12,
    replies: 8,
    views: 145,
    createdAt: '2024-01-10',
  },
  {
    id: '5',
    title: 'Análisis técnico básico: Patrones que funcionan',
    content: 'Voy a explicar los patrones de gráficos más comunes y cómo identificarlos...',
    author: { username: 'ChartMaster', rank: 'Oro 3' },
    category: 'strategies',
    likes: 178,
    replies: 56,
    views: 2100,
    createdAt: '2024-01-07',
  },
  {
    id: '6',
    title: 'Presentación: Soy nuevo en la comunidad!',
    content: 'Hola! Me llamo Juan y acabo de descubrir SmartStocks. Estoy emocionado...',
    author: { username: 'JuanInversor', rank: 'Bronce 3' },
    category: 'general',
    likes: 45,
    replies: 32,
    views: 320,
    createdAt: '2024-01-09',
  },
]

const categories = [
  { id: 'all', label: 'Todos', icon: MessageSquare },
  { id: 'general', label: 'General', icon: MessageCircle },
  { id: 'strategies', label: 'Estrategias', icon: Lightbulb },
  { id: 'help', label: 'Ayuda', icon: HelpCircle },
  { id: 'news', label: 'Noticias', icon: TrendingUp },
]

const categoryColors: Record<string, string> = {
  general: 'bg-gray-100 text-gray-700',
  strategies: 'bg-purple-100 text-purple-700',
  help: 'bg-blue-100 text-blue-700',
  news: 'bg-green-100 text-green-700',
}

const categoryLabels: Record<string, string> = {
  general: 'General',
  strategies: 'Estrategias',
  help: 'Ayuda',
  news: 'Noticias',
}

function PostCard({ post }: { post: ForumPost }) {
  return (
    <Card className="hover:shadow-md transition-shadow cursor-pointer">
      <CardContent className="p-4">
        <div className="flex gap-4">
          <Avatar className="h-10 w-10">
            <AvatarImage src={post.author.avatar} />
            <AvatarFallback>
              {post.author.username.charAt(0).toUpperCase()}
            </AvatarFallback>
          </Avatar>

          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-2">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 flex-wrap">
                  {post.isPinned && (
                    <Pin className="w-4 h-4 text-primary flex-shrink-0" />
                  )}
                  {post.isHot && (
                    <Flame className="w-4 h-4 text-orange-500 flex-shrink-0" />
                  )}
                  <h3 className="font-semibold truncate">{post.title}</h3>
                </div>
                <p className="text-sm text-muted-foreground line-clamp-1 mt-1">
                  {post.content}
                </p>
              </div>
              <Badge className={cn('flex-shrink-0', categoryColors[post.category])}>
                {categoryLabels[post.category]}
              </Badge>
            </div>

            <div className="flex items-center gap-4 mt-3 text-sm text-muted-foreground">
              <span className="flex items-center gap-1">
                <span className="font-medium text-foreground">{post.author.username}</span>
                <span className="text-xs">({post.author.rank})</span>
              </span>
              <span className="flex items-center gap-1">
                <ThumbsUp className="w-4 h-4" />
                {post.likes}
              </span>
              <span className="flex items-center gap-1">
                <MessageCircle className="w-4 h-4" />
                {post.replies}
              </span>
              <span className="flex items-center gap-1">
                <Eye className="w-4 h-4" />
                {post.views}
              </span>
              <span className="flex items-center gap-1 ml-auto">
                <Clock className="w-4 h-4" />
                {post.createdAt}
              </span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

const STORAGE_KEY = 'smartstocks-forum-posts'

export default function ForumPage() {
  const { user } = useAuthStore()
  const [selectedCategory, setSelectedCategory] = useState('all')
  const [searchQuery, setSearchQuery] = useState('')
  const [posts, setPosts] = useState<ForumPost[]>(() => {
    // Inicializar desde localStorage si existe
    if (typeof window !== 'undefined') {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (saved) {
        try {
          return JSON.parse(saved)
        } catch {
          return mockPosts
        }
      }
    }
    return mockPosts
  })
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [newPost, setNewPost] = useState({
    title: '',
    content: '',
    category: 'general' as ForumPost['category'],
  })

  // Guardar posts en localStorage cuando cambien
  useEffect(() => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(posts))
  }, [posts])

  const handleCreatePost = () => {
    if (!newPost.title.trim() || !newPost.content.trim()) return

    const post: ForumPost = {
      id: Date.now().toString(),
      title: newPost.title.trim(),
      content: newPost.content.trim(),
      category: newPost.category,
      author: {
        username: user?.username || 'Usuario',
        rank: 'Bronce 1',
      },
      likes: 0,
      replies: 0,
      views: 1,
      createdAt: new Date().toISOString().split('T')[0],
    }

    // Agregar el post primero
    setPosts((prevPosts) => {
      const newPosts = [post, ...prevPosts]
      console.log('Posts actualizados:', newPosts.length)
      return newPosts
    })

    // Luego actualizar el resto
    setSelectedCategory('all')
    setIsDialogOpen(false)
    setNewPost({ title: '', content: '', category: 'general' })

    toast.success('Post publicado correctamente')
  }

  const filteredPosts = posts.filter((post) => {
    const matchesCategory = selectedCategory === 'all' || post.category === selectedCategory
    const matchesSearch = post.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      post.content.toLowerCase().includes(searchQuery.toLowerCase())
    return matchesCategory && matchesSearch
  })

  // Sort: pinned first, then by date (newest first)
  const sortedPosts = [...filteredPosts].sort((a, b) => {
    if (a.isPinned && !b.isPinned) return -1
    if (!a.isPinned && b.isPinned) return 1
    // Sort by date descending (newest first)
    return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  })

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold text-foreground flex items-center gap-3">
            <MessageSquare className="w-8 h-8 text-primary" />
            Foro de la Comunidad
          </h1>
          <p className="text-muted-foreground mt-1">
            Comparte estrategias, haz preguntas y conecta con otros traders
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="w-4 h-4 mr-2" />
              Nuevo Post
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>Crear nuevo post</DialogTitle>
              <DialogDescription>
                Comparte tus ideas, estrategias o preguntas con la comunidad.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="title">Título</Label>
                <Input
                  id="title"
                  placeholder="Escribe un título descriptivo..."
                  value={newPost.title}
                  onChange={(e) => setNewPost({ ...newPost, title: e.target.value })}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="category">Categoría</Label>
                <Select
                  value={newPost.category}
                  onValueChange={(value: ForumPost['category']) =>
                    setNewPost({ ...newPost, category: value })
                  }
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Selecciona una categoría" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="general">
                      <div className="flex items-center gap-2">
                        <MessageCircle className="w-4 h-4" />
                        General
                      </div>
                    </SelectItem>
                    <SelectItem value="strategies">
                      <div className="flex items-center gap-2">
                        <Lightbulb className="w-4 h-4" />
                        Estrategias
                      </div>
                    </SelectItem>
                    <SelectItem value="help">
                      <div className="flex items-center gap-2">
                        <HelpCircle className="w-4 h-4" />
                        Ayuda
                      </div>
                    </SelectItem>
                    <SelectItem value="news">
                      <div className="flex items-center gap-2">
                        <TrendingUp className="w-4 h-4" />
                        Noticias
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="space-y-2">
                <Label htmlFor="content">Contenido</Label>
                <Textarea
                  id="content"
                  placeholder="Escribe el contenido de tu post..."
                  rows={5}
                  value={newPost.content}
                  onChange={(e) => setNewPost({ ...newPost, content: e.target.value })}
                />
              </div>
            </div>
            <DialogFooter>
              <Button variant="outline" onClick={() => setIsDialogOpen(false)}>
                Cancelar
              </Button>
              <Button
                onClick={handleCreatePost}
                disabled={!newPost.title.trim() || !newPost.content.trim()}
              >
                Publicar
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>

      {/* Search */}
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
        <Input
          placeholder="Buscar en el foro..."
          className="pl-10"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4 text-center">
            <MessageSquare className="w-6 h-6 text-primary mx-auto mb-2" />
            <p className="text-2xl font-bold">1,234</p>
            <p className="text-sm text-muted-foreground">Posts Totales</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 text-center">
            <MessageCircle className="w-6 h-6 text-blue-500 mx-auto mb-2" />
            <p className="text-2xl font-bold">5,678</p>
            <p className="text-sm text-muted-foreground">Respuestas</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 text-center">
            <Users className="w-6 h-6 text-green-500 mx-auto mb-2" />
            <p className="text-2xl font-bold">890</p>
            <p className="text-sm text-muted-foreground">Miembros Activos</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4 text-center">
            <Flame className="w-6 h-6 text-orange-500 mx-auto mb-2" />
            <p className="text-2xl font-bold">45</p>
            <p className="text-sm text-muted-foreground">Posts Hoy</p>
          </CardContent>
        </Card>
      </div>

      {/* Categories */}
      <Tabs value={selectedCategory} onValueChange={setSelectedCategory}>
        <TabsList className="w-full justify-start overflow-x-auto">
          {categories.map((category) => (
            <TabsTrigger
              key={category.id}
              value={category.id}
              className="flex items-center gap-2"
            >
              <category.icon className="w-4 h-4" />
              {category.label}
            </TabsTrigger>
          ))}
        </TabsList>

        <TabsContent value={selectedCategory} className="mt-6">
          <div className="space-y-4">
            {sortedPosts.length > 0 ? (
              sortedPosts.map((post) => (
                <PostCard key={post.id} post={post} />
              ))
            ) : (
              <div className="text-center py-12">
                <MessageSquare className="w-12 h-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-medium">No hay posts</h3>
                <p className="text-muted-foreground">
                  {searchQuery
                    ? 'No se encontraron posts con tu búsqueda'
                    : 'Sé el primero en crear un post en esta categoría'}
                </p>
              </div>
            )}
          </div>
        </TabsContent>
      </Tabs>

      {/* Coming Soon Notice */}
      <Card className="bg-gradient-to-r from-blue-50 to-purple-50 border-blue-200">
        <CardContent className="py-6">
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
              <MessageSquare className="w-6 h-6 text-primary" />
            </div>
            <div>
              <h3 className="font-semibold">Foro en Desarrollo</h3>
              <p className="text-sm text-muted-foreground">
                Esta es una vista previa del foro. La funcionalidad completa estará disponible pronto.
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
