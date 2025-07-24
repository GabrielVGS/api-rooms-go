import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { ArrowLeft, Users, Calendar, Plus, Edit, Trash2, UserPlus, UserMinus } from 'lucide-react';
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

  const roomId = id ? parseInt(id, 10) : 0;

  const fetchRoomDetails = async () => {
    if (!roomId) return;
    
    try {
      setLoading(true);
      const roomData = await roomApi.getRoomById(roomId);
      setRoom(roomData);
      
      // Check if current user is a member
      const currentUser = JSON.parse(localStorage.getItem('user') || '{}');
      const isMember = roomData.members?.some(member => member.user_id === currentUser.id) || false;
      setIsUserMember(isMember);
      
      setError(null);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to fetch room details');
    } finally {
      setLoading(false);
    }
  };

  const fetchNotes = async () => {
    if (!roomId || !isUserMember) return;
    
    try {
      const notesData = await notesApi.getNotesByRoom(roomId);
      setNotes(notesData);
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
      fetchRoomDetails(); // Refresh room data to show updated members
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to join room');
    }
  };

  const handleLeaveRoom = async () => {
    if (!confirm('Are you sure you want to leave this room?')) return;
    
    try {
      await roomApi.leaveRoom(roomId);
      setIsUserMember(false);
      setNotes([]); // Clear notes when leaving
      fetchRoomDetails(); // Refresh room data
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to leave room');
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
    if (!confirm('Are you sure you want to delete this note?')) return;

    try {
      await notesApi.deleteNote(noteId);
      setNotes(notes.filter(note => note.id !== noteId));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete note');
    }
  };

  const handleNoteSaved = () => {
    fetchNotes();
    setNoteDialogOpen(false);
    setEditingNote(null);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-muted-foreground">Loading room details...</div>
      </div>
    );
  }

  if (!room) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-muted-foreground">Room not found</div>
      </div>
    );
  }

  const currentUser = JSON.parse(localStorage.getItem('user') || '{}');

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <Button
            variant="outline"
            size="sm"
            onClick={() => navigate('/rooms')}
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Rooms
          </Button>
          <div>
            <h1 className="text-3xl font-bold">{room.name}</h1>
            <p className="text-muted-foreground">{room.subject}</p>
          </div>
        </div>
        
        <div className="flex space-x-2">
          {isUserMember ? (
            <Button variant="outline" onClick={handleLeaveRoom}>
              <UserMinus className="mr-2 h-4 w-4" />
              Leave Room
            </Button>
          ) : (
            <Button onClick={handleJoinRoom}>
              <UserPlus className="mr-2 h-4 w-4" />
              Join Room
            </Button>
          )}
        </div>
      </div>

      {error && (
        <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Room Information */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Calendar className="mr-2 h-5 w-5" />
              Room Details
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <label className="text-sm font-medium text-muted-foreground">Description</label>
              <p className="text-sm">{room.description || 'No description provided'}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-muted-foreground">Capacity</label>
              <p className="text-sm">{room.capacity} people</p>
            </div>
            <div>
              <label className="text-sm font-medium text-muted-foreground">Created</label>
              <p className="text-sm">{new Date(room.created_at).toLocaleDateString()}</p>
            </div>
          </CardContent>
        </Card>

        {/* Members */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Users className="mr-2 h-5 w-5" />
              Members ({room.members?.length || 0})
            </CardTitle>
          </CardHeader>
          <CardContent>
            {room.members && room.members.length > 0 ? (
              <div className="space-y-2">
                {room.members.map((member) => (
                  <div key={member.user_id} className="flex items-center justify-between p-2 bg-muted/50 rounded">
                    <div>
                      <p className="text-sm font-medium">{member.user_name}</p>
                      <p className="text-xs text-muted-foreground">{member.user_email}</p>
                    </div>
                    <span className="text-xs bg-primary/10 text-primary px-2 py-1 rounded">
                      {member.role}
                    </span>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">No members yet</p>
            )}
          </CardContent>
        </Card>

        {/* Notes Section */}
        <Card className="md:col-span-1">
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="flex items-center">
                Notes ({notes.length})
              </CardTitle>
              {isUserMember && (
                <Button size="sm" onClick={handleCreateNote}>
                  <Plus className="mr-2 h-4 w-4" />
                  Add Note
                </Button>
              )}
            </div>
          </CardHeader>
          <CardContent>
            {!isUserMember ? (
              <p className="text-sm text-muted-foreground">
                Join this room to view and create notes
              </p>
            ) : notes.length > 0 ? (
              <div className="space-y-3">
                {notes.map((note) => (
                  <div key={note.id} className="border rounded-lg p-3">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <h4 className="font-medium text-sm">{note.title}</h4>
                        <p className="text-xs text-muted-foreground mt-1">
                          {note.content.substring(0, 100)}
                          {note.content.length > 100 ? '...' : ''}
                        </p>
                        <p className="text-xs text-muted-foreground mt-2">
                          By {note.user_name} • {new Date(note.created_at).toLocaleDateString()}
                        </p>
                      </div>
                      {note.user_id === currentUser.id && (
                        <div className="flex space-x-1 ml-2">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleEditNote(note)}
                          >
                            <Edit className="h-3 w-3" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleDeleteNote(note.id)}
                          >
                            <Trash2 className="h-3 w-3" />
                          </Button>
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-muted-foreground">
                No notes yet. Create the first note for this room!
              </p>
            )}
          </CardContent>
        </Card>
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