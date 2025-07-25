import React, { useState, useEffect } from 'react';
import { Plus, Edit, Trash2, Calendar, Clock } from 'lucide-react';
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
import { reservationApi, roomApi } from '@/services/api';
import type { Reservation, Room } from '@/types/api';
import { ReservationDialog } from './ReservationDialog';

export const ReservationList: React.FC = () => {
  const [reservations, setReservations] = useState<Reservation[]>([]);
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);
  const [editingReservation, setEditingReservation] = useState<Reservation | null>(null);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [reservationsData, roomsData] = await Promise.all([
        reservationApi.getReservations(),
        roomApi.getRooms(),
      ]);
      setReservations(reservationsData);
      setRooms(roomsData);
      setError(null);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Falha ao carregar dados');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, []);

  const handleDelete = async (id: number) => {
    if (!confirm('Tem certeza que deseja cancelar esta reserva?')) return;

    try {
      await reservationApi.deleteReservation(id);
      setReservations(reservations.filter(reservation => reservation.id !== id));
    } catch (err: any) {
      setError(err.response?.data?.message || 'Falha ao cancelar reserva');
    }
  };

  const handleEdit = (reservation: Reservation) => {
    setEditingReservation(reservation);
    setDialogOpen(true);
  };

  const handleCreate = () => {
    setEditingReservation(null);
    setDialogOpen(true);
  };

  const handleReservationSaved = () => {
    fetchData();
    setDialogOpen(false);
    setEditingReservation(null);
  };

  const getRoomName = (roomId: number) => {
    const room = rooms.find(r => r.id === roomId);
    return room?.name || 'Sala Desconhecida';
  };


  const formatTime = (dateTime: string) => {
    return new Date(dateTime).toLocaleTimeString([], { 
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };

  const formatDate = (dateTime: string) => {
    return new Date(dateTime).toLocaleDateString();
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-muted-foreground">Carregando reservas...</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Reservas</h1>
          <p className="text-muted-foreground">Gerencie reservas de salas e agendamentos</p>
        </div>
        <Button onClick={handleCreate}>
          <Plus className="mr-2 h-4 w-4" />
          Nova Reserva
        </Button>
      </div>

      {error && (
        <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </div>
      )}

      <Card>
        <CardHeader>
          <CardTitle>Suas Reservas</CardTitle>
        </CardHeader>
        <CardContent>
          {reservations.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              Nenhuma reserva encontrada. Crie sua primeira reserva para começar.
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Sala</TableHead>
                  <TableHead>Data</TableHead>
                  <TableHead>Horário</TableHead>
                  <TableHead>Duração</TableHead>
                  <TableHead className="text-right">Ações</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {reservations.map((reservation) => (
                  <TableRow key={reservation.id}>
                    <TableCell className="font-medium">
                      {getRoomName(reservation.roomId)}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center">
                        <Calendar className="mr-2 h-4 w-4 text-muted-foreground" />
                        {formatDate(reservation.startTime)}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center">
                        <Clock className="mr-2 h-4 w-4 text-muted-foreground" />
                        {formatTime(reservation.startTime)} - {formatTime(reservation.endTime)}
                      </div>
                    </TableCell>
                    <TableCell>
                      {Math.round(
                        (new Date(reservation.endTime).getTime() - 
                         new Date(reservation.startTime).getTime()) / 
                        (1000 * 60)
                      )} min
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end space-x-2">
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleEdit(reservation)}
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => handleDelete(reservation.id)}
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

      <ReservationDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        reservation={editingReservation}
        rooms={rooms}
        onSave={handleReservationSaved}
      />
    </div>
  );
};