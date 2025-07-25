import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  ArrowLeft, 
  Users, 
  Calendar, 
  Plus, 
  Edit, 
  Trash2, 
  UserPlus, 
  UserMinus,
  MessageSquare,
  Clock,
  User,
  MoreVertical,
  Search,
  Filter,
  FileText
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { roomApi, notesApi } from '@/services/api';
import type { Room, Note } from '@/types/api';
import { NoteDialog } from './NoteDialog';

export const RoomDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [room, setRoom] = useState<Room | null>(null);
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [noteDialogOpen, setNoteDialogOpen] = useState(false);
  const [editingNote, setEditingNote] = useState<Note | null>(null);
  const [isUserMember, setIsUserMember] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [expandedNote, setExpandedNote] = useState<number | null>(null);

  const roomId = id ? parseInt(id, 10) : 0;

  const fetchRoomDetails = async () => {
    if (!roomId) return;
    
    try {
      setLoading(true);
      const roomData = await roomApi.getRoomById(roomId);
      setRoom(roomData);
      
      // Check if current user is a member
      const currentUser = JSON.parse(localStorage.getItem('auth_user') || '{}');
      const isMember = roomData.members?.some(member => member.user_id === currentUser.id) || false;
      setIsUserMember(isMember);
      
      setError(null);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Falha ao carregar detalhes da sala');
    } finally {
      setLoading(false);
    }
  };

  const fetchNotes = async () => {
    if (!roomId || !isUserMember) return;
    
    try {
      const notesData = await notesApi.getNotesByRoom(roomId);
      setNotes(notesData.reverse());
    } catch (err: any) {
      console.error('Failed to fetch notes:', err);
    }
  };

  useEffect(() => {
    fetchRoomDetails();
  }, [roomId]);

  useEffect(() => {
    if (isUserMember) {
      fetchNotes();
    }
  }, [isUserMember, roomId]);

  const handleJoinRoom = async () => {
    try {
      await roomApi.joinRoom(roomId);
      setIsUserMember(true);
      fetchRoomDetails();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Sala cheia');
    }
  };

  const handleLeaveRoom = async () => {
    if (!confirm('Tem certeza que deseja sair desta sala?')) return;
    
    try {
      await roomApi.leaveRoom(roomId);
      setIsUserMember(false);
      setNotes([]); 
      fetchRoomDetails(); 
    } catch (err: any) {
      setError(err.response?.data?.message || 'Falha ao sair da sala');
    }
  };

  const handleCreateNote = () => {
    setEditingNote(null);
    setNoteDialogOpen(true);
  };

  const handleEditNote = (note: Note) => {
    setEditingNote(note);
    setNoteDialogOpen(true);
  };

  const handleDeleteNote = async (noteId: number) => {
    if (!confirm('Tem certeza que deseja excluir esta nota?')) return;

    try {
      await notesApi.deleteNote(noteId);
      setNotes(notes.filter(note => note.id !== noteId));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Falha ao excluir nota');
    }
  };

  const handleNoteSaved = () => {
    fetchNotes();
    setNoteDialogOpen(false);
    setEditingNote(null);
  };

  // Filter notes based on search query
  const filteredNotes = notes.filter(note => 
    note.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
    note.content.toLowerCase().includes(searchQuery.toLowerCase()) ||
    note.user_name?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Carregando detalhes da sala...</p>
        </div>
      </div>
    );
  }

  if (!room) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <FileText className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
          <h2 className="text-xl font-semibold mb-2">Sala não encontrada</h2>
          <p className="text-muted-foreground mb-4">A sala que você está procurando não existe ou foi removida.</p>
          <Button onClick={() => navigate('/rooms')} variant="outline">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Voltar para Salas
          </Button>
        </div>
      </div>
    );
  }

  const currentUser = JSON.parse(localStorage.getItem('auth_user') || '{}');

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted/20">
      {/* Enhanced Header */}
      <div className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky top-0 z-50">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => navigate('/rooms')}
                className="hover:bg-muted/50"
              >
                <ArrowLeft className="mr-2 h-4 w-4" />
                Voltar para Salas
              </Button>
              <div className="border-l pl-4">
                <div className="flex items-center space-x-3">
                  <div className="w-12 h-12 bg-gradient-to-br from-primary/20 to-primary/5 rounded-xl flex items-center justify-center">
                    <MessageSquare className="h-6 w-6 text-primary" />
                  </div>
                  <div>
                    <h1 className="text-2xl md:text-3xl font-bold tracking-tight">{room.name}</h1>
                    <p className="text-muted-foreground flex items-center mt-1">
                      <span className="bg-secondary/60 px-2 py-1 rounded-md text-xs font-medium mr-2">
                        {room.subject}
                      </span>
                      <Users className="h-3 w-3 mr-1" />
                      {room.members?.length || 0} membros
                      <Calendar className="h-3 w-3 ml-3 mr-1" />
                      {room.capacity} capacidade
                    </p>
                  </div>
                </div>
              </div>
            </div>
            
            <div className="flex items-center space-x-3">
              {isUserMember ? (
                <>
                  <Button variant="outline" onClick={handleCreateNote} className="hidden sm:flex">
                    <Plus className="mr-2 h-4 w-4" />
                    Adicionar Nota
                  </Button>
                  <Button variant="outline" onClick={handleLeaveRoom} className="text-destructive hover:text-destructive">
                    <UserMinus className="mr-2 h-4 w-4" />
                    <span className="hidden sm:inline">Sair da Sala</span>
                  </Button>
                </>
              ) : (
                <Button onClick={handleJoinRoom} className="bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70">
                  <UserPlus className="mr-2 h-4 w-4" />
                  Entrar na Sala
                </Button>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-6 py-8">
        {error && (
          <div className="mb-6 rounded-lg bg-destructive/10 border border-destructive/20 p-4 text-sm text-destructive flex items-center">
            <div className="w-2 h-2 bg-destructive rounded-full mr-3"></div>
            {error}
          </div>
        )}

        {/* Responsive Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Sidebar with Room Info and Members */}
          <div className="lg:col-span-1 space-y-6">
            {/* Room Information */}
            <Card className="backdrop-blur-sm bg-background/95 border-muted/50">
              <CardHeader className="pb-3">
                <CardTitle className="flex items-center text-base">
                  <Calendar className="mr-2 h-4 w-4 text-primary" />
                  Detalhes da Sala
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4 text-sm">
                <div className="space-y-2">
                  <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Descrição</label>
                  <p className="text-sm leading-relaxed">
                    {room.description || 'Nenhuma descrição fornecida'}
                  </p>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Capacidade</label>
                    <p className="text-sm font-medium">{room.capacity} pessoas</p>
                  </div>
                  <div>
                    <label className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Criada em</label>
                    <p className="text-sm font-medium">{new Date(room.created_at).toLocaleDateString()}</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Members */}
            <Card className="backdrop-blur-sm bg-background/95 border-muted/50">
              <CardHeader className="pb-3">
                <CardTitle className="flex items-center justify-between text-base">
                  <span className="flex items-center">
                    <Users className="mr-2 h-4 w-4 text-primary" />
                    Membros
                  </span>
                  <span className="text-xs bg-muted px-2 py-1 rounded-full font-medium">
                    {room.members?.length || 0}
                  </span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                {room.members && room.members.length > 0 ? (
                  <div className="space-y-3 max-h-80 overflow-y-auto">
                    {room.members.map((member) => (
                      <div key={member.user_id} className="flex items-center space-x-3 p-3 bg-muted/30 rounded-lg hover:bg-muted/50 transition-colors">
                        <div className="w-8 h-8 bg-gradient-to-br from-primary/20 to-primary/5 rounded-full flex items-center justify-center">
                          <User className="h-4 w-4 text-primary" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="text-sm font-medium truncate">{member.user_name}</p>
                          <p className="text-xs text-muted-foreground truncate">{member.user_email}</p>
                        </div>
                        <span className="text-xs bg-primary/10 text-primary px-2 py-1 rounded-md font-medium">
                          {member.role}
                        </span>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <Users className="h-8 w-8 text-muted-foreground mx-auto mb-2" />
                    <p className="text-sm text-muted-foreground">Nenhum membro ainda</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Main Content - Notes */}
          <div className="lg:col-span-3">
            <Card className="backdrop-blur-sm bg-background/95 border-muted/50 min-h-[600px]">
              <CardHeader className="pb-4 border-b border-muted/50">
                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                  <div className="flex items-center space-x-2">
                    <MessageSquare className="h-5 w-5 text-primary" />
                    <CardTitle className="text-lg">
                      Notas
                      <span className="ml-2 text-xs bg-muted px-2 py-1 rounded-full font-medium">
                        {filteredNotes.length}
                      </span>
                    </CardTitle>
                  </div>
                  
                  <div className="flex flex-col sm:flex-row gap-2 sm:items-center">
                    {isUserMember && (
                      <>
                        {/* Search Bar */}
                        <div className="relative">
                          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                          <input
                            type="text"
                            placeholder="Buscar notas..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="pl-10 pr-4 py-2 text-sm border border-muted/50 rounded-lg bg-background/50 focus:ring-2 focus:ring-primary/20 focus:border-primary/30 transition-all w-full sm:w-64"
                          />
                        </div>
                        
                        {/* Add Note Button - Mobile */}
                        <Button onClick={handleCreateNote} className="sm:hidden w-full">
                          <Plus className="mr-2 h-4 w-4" />
                          Adicionar Nota
                        </Button>
                      </>
                    )}
                  </div>
                </div>
              </CardHeader>
              
              <CardContent className="p-0">
                {!isUserMember ? (
                  <div className="text-center py-16">
                    <UserPlus className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                    <h3 className="text-lg font-semibold mb-2">Entre na sala para acessar as notas</h3>
                    <p className="text-muted-foreground mb-6 max-w-sm mx-auto">
                      Torne-se um membro para visualizar, criar e colaborar em notas com outros membros da sala.
                    </p>
                    <Button onClick={handleJoinRoom} className="bg-gradient-to-r from-primary to-primary/80">
                      <UserPlus className="mr-2 h-4 w-4" />
                      Entrar na Sala
                    </Button>
                  </div>
                ) : filteredNotes.length > 0 ? (
                  <div className="divide-y divide-muted/30">
                    {filteredNotes.map((note) => (
                      <div key={note.id} className="p-6 hover:bg-muted/20 transition-colors group">
                        <div className="flex items-start justify-between">
                          <div className="flex-1 min-w-0">
                            <div className="flex items-center space-x-2 mb-2">
                              <h4 className="font-semibold text-base text-foreground group-hover:text-primary transition-colors">
                                {note.title}
                              </h4>
                              {note.user_id === currentUser.id && (
                                <span className="text-xs bg-primary/10 text-primary px-2 py-1 rounded-md font-medium">
                                  Sua nota
                                </span>
                              )}
                            </div>
                            
                            <div className="prose prose-sm max-w-none mb-3">
                              <p className="text-muted-foreground leading-relaxed">
                                {expandedNote === note.id ? note.content : (
                                  <>
                                    {note.content.substring(0, 200)}
                                    {note.content.length > 200 && (
                                      <button
                                        onClick={() => setExpandedNote(note.id)}
                                        className="text-primary hover:underline ml-1 font-medium"
                                      >
                                        Ler mais
                                      </button>
                                    )}
                                  </>
                                )}
                                {expandedNote === note.id && note.content.length > 200 && (
                                  <button
                                    onClick={() => setExpandedNote(null)}
                                    className="text-primary hover:underline ml-2 font-medium"
                                  >
                                    Mostrar menos
                                  </button>
                                )}
                              </p>
                            </div>
                            
                            <div className="flex items-center text-xs text-muted-foreground space-x-4">
                              <div className="flex items-center">
                                <User className="h-3 w-3 mr-1" />
                                {note.user_name}
                              </div>
                              <div className="flex items-center">
                                <Clock className="h-3 w-3 mr-1" />
                                {new Date(note.created_at).toLocaleDateString('en-US', { 
                                  year: 'numeric',
                                  month: 'short',
                                  day: 'numeric',
                                  hour: '2-digit',
                                  minute: '2-digit'
                                })}
                              </div>
                            </div>
                          </div>
                          
                          {note.user_id === currentUser.id && (
                            <div className="flex items-center space-x-1 ml-4 opacity-0 group-hover:opacity-100 transition-opacity">
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleEditNote(note)}
                                className="h-8 w-8 p-0 hover:bg-primary/10 hover:text-primary"
                              >
                                <Edit className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleDeleteNote(note.id)}
                                className="h-8 w-8 p-0 hover:bg-destructive/10 hover:text-destructive"
                              >
                                <Trash2 className="h-4 w-4" />
                              </Button>
                            </div>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-16">
                    {searchQuery ? (
                      <>
                        <Search className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                        <h3 className="text-lg font-semibold mb-2">Nenhuma nota encontrada</h3>
                        <p className="text-muted-foreground mb-6">
                          Nenhuma nota corresponde à sua busca "{searchQuery}"
                        </p>
                        <Button
                          variant="outline"
                          onClick={() => setSearchQuery('')}
                        >
                          Limpar busca
                        </Button>
                      </>
                    ) : (
                      <>
                        <FileText className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                        <h3 className="text-lg font-semibold mb-2">Nenhuma nota ainda</h3>
                        <p className="text-muted-foreground mb-6 max-w-sm mx-auto">
                          Seja o primeiro a criar uma nota nesta sala. Compartilhe seus pensamentos, ideias ou informações importantes com outros membros.
                        </p>
                        <Button onClick={handleCreateNote} className="bg-gradient-to-r from-primary to-primary/80">
                          <Plus className="mr-2 h-4 w-4" />
                          Criar Primeira Nota
                        </Button>
                      </>
                    )}
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      <NoteDialog
        open={noteDialogOpen}
        onOpenChange={setNoteDialogOpen}
        note={editingNote}
        roomId={roomId}
        onSave={handleNoteSaved}
      />
    </div>
  );
};