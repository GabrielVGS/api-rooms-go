import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { roomApi } from '@/services/api';
import type { Room } from '@/types/api';
import { RoomDialog } from './RoomDialog';




export const RoomList: React.FC = () => {
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingRoom, setEditingRoom] = useState<Room | null>(null);

  const fetchRooms = async () => {
    try {
      setLoading(true);
      const data = await roomApi.getRooms();
      setRooms(data);
      setError(null);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to fetch rooms');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRooms();
  }, []);

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this room?')) return;

    try {
      await roomApi.deleteRoom(id);
      setRooms(rooms.filter(room => room.id !== id));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to delete room');
    }
  };

  const handleEdit = (room: Room) => {
    setEditingRoom(room);
    setDialogOpen(true);
  };

  const handleCreate = () => {
    setEditingRoom(null);
    setDialogOpen(true);
  };

  const handleRoomSaved = () => {
    fetchRooms();
    setDialogOpen(false);
    setEditingRoom(null);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-muted-foreground">Loading rooms...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Rooms</h1>
          <p className="text-muted-foreground">Manage meeting rooms and spaces</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="mr-2 h-4 w-4" />
          Add Room
        </Button>
      </div>

      {error && (
        <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </div>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Available Rooms</CardTitle>
        </CardHeader>
        <CardContent>
          {rooms.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              No rooms found. Create your first room to get started.
            </div>
          ) : (

            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Capacity</TableHead>
                  <TableHead>Assunto</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rooms.map((room) => (
                  <TableRow key={room.id}>
                    <TableCell className="font-medium">{room.name}</TableCell>
                    <TableCell>{room.capacity} people</TableCell>
                    <TableCell>{room.subject}</TableCell>
                    <TableCell className="max-w-xs truncate">
                      {room.description || 'No description'}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end space-x-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleEdit(room)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleDelete(room.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      <RoomDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        room={editingRoom}
        onSave={handleRoomSaved}
      />
    </div>
  );
};