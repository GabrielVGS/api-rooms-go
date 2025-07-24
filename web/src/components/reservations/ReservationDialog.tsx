import React, { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { reservationApi } from '@/services/api';
import type { Reservation, Room } from '@/types/api';

const reservationSchema = z.object({
  roomId: z.number().min(1, 'Please select a room'),
  startTime: z.string().min(1, 'Start time is required'),
  endTime: z.string().min(1, 'End time is required'),
}).refine((data) => {
  const start = new Date(data.startTime);
  const end = new Date(data.endTime);
  return end > start;
}, {
  message: "End time must be after start time",
  path: ["endTime"],
});

type ReservationFormData = z.infer<typeof reservationSchema>;

interface ReservationDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  reservation: Reservation | null;
  rooms: Room[];
  onSave: () => void;
}

export const ReservationDialog: React.FC<ReservationDialogProps> = ({
  open,
  onOpenChange,
  reservation,
  rooms,
  onSave,
}) => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const isEditing = !!reservation;

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
    setValue,
  } = useForm<ReservationFormData>({
    resolver: zodResolver(reservationSchema),
  });

  useEffect(() => {
    if (open) {
      if (reservation) {
        setValue('roomId', reservation.roomId);
        setValue('startTime', new Date(reservation.startTime).toISOString().slice(0, 16));
        setValue('endTime', new Date(reservation.endTime).toISOString().slice(0, 16));
      } else {
        reset();
        const now = new Date();
        const oneHourLater = new Date(now.getTime() + 60 * 60 * 1000);
        setValue('startTime', now.toISOString().slice(0, 16));
        setValue('endTime', oneHourLater.toISOString().slice(0, 16));
      }
      setError(null);
    }
  }, [open, reservation, setValue, reset]);

  const onSubmit = async (data: ReservationFormData) => {
    setIsLoading(true);
    setError(null);

    try {
      const reservationData = {
        ...data,
        startTime: new Date(data.startTime).toISOString(),
        endTime: new Date(data.endTime).toISOString(),
      };

      if (isEditing && reservation) {
        await reservationApi.updateReservation(reservation.id, reservationData);
      } else {
        await reservationApi.createReservation({
          ...reservationData,
          userId: 0, // This should be set by the backend based on the auth token
        });
      }
      onSave();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to save reservation');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Edit Reservation' : 'Create New Reservation'}
          </DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Update the reservation details below.'
              : 'Book a room for your meeting or event.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          {error && (
            <div className="rounded-md bg-destructive/10 p-3 text-sm text-destructive">
              {error}
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="roomId">Room</Label>
            <select
              id="roomId"
              {...register('roomId', { valueAsNumber: true })}
              className={`flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 ${
                errors.roomId ? 'border-destructive' : ''
              }`}
            >
              <option value="">Select a room</option>
              {rooms.map((room) => (
                <option key={room.id} value={room.id}>
                  {room.name} (Capacity: {room.capacity})
                </option>
              ))}
            </select>
            {errors.roomId && (
              <p className="text-sm text-destructive">{errors.roomId.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="startTime">Start Time</Label>
            <Input
              id="startTime"
              type="datetime-local"
              {...register('startTime')}
              className={errors.startTime ? 'border-destructive' : ''}
            />
            {errors.startTime && (
              <p className="text-sm text-destructive">{errors.startTime.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="endTime">End Time</Label>
            <Input
              id="endTime"
              type="datetime-local"
              {...register('endTime')}
              className={errors.endTime ? 'border-destructive' : ''}
            />
            {errors.endTime && (
              <p className="text-sm text-destructive">{errors.endTime.message}</p>
            )}
          </div>

          <div className="flex justify-end space-x-2 pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading
                ? isEditing
                  ? 'Updating...'
                  : 'Creating...'
                : isEditing
                ? 'Update Reservation'
                : 'Create Reservation'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};