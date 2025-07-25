import React, { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { FileText, Loader2 } from 'lucide-react';
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
import { notesApi } from '@/services/api';
import type { Note, CreateNoteRequest, UpdateNoteRequest } from '@/types/api';

const noteSchema = z.object({
  title: z.string().min(1, 'Título é obrigatório').max(100, 'Título deve ter menos de 100 caracteres'),
  content: z.string().min(1, 'Conteúdo é obrigatório').max(5000, 'Conteúdo deve ter menos de 5000 caracteres'),
});

type NoteFormData = z.infer<typeof noteSchema>;

interface NoteDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  note: Note | null;
  roomId: number;
  onSave: () => void;
}

export const NoteDialog: React.FC<NoteDialogProps> = ({
  open,
  onOpenChange,
  note,
  roomId,
  onSave,
}) => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const isEditing = !!note;

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
    setValue,
  } = useForm<NoteFormData>({
    resolver: zodResolver(noteSchema),
  });

  useEffect(() => {
    if (open) {
      if (note) {
        setValue('title', note.title);
        setValue('content', note.content);
      } else {
        reset({
          title: '',
          content: ''
        });
      }
      setError(null);
    }
  }, [open, note, setValue, reset]);

  const onSubmit = async (data: NoteFormData) => {
    setIsLoading(true);
    setError(null);

    try {
      if (isEditing && note) {
        const updateRequest: UpdateNoteRequest = {
          title: data.title,
          content: data.content,
        };
        await notesApi.updateNote(note.id, updateRequest);
      } else {
        const createRequest: CreateNoteRequest = {
          room_id: roomId,
          title: data.title,
          content: data.content,
        };
        await notesApi.createNote(createRequest);
      }
      onSave();
    } catch (err: any) {
      setError(err.response?.data?.message || 'Falha ao salvar nota');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
        <DialogHeader className="pb-4 border-b">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-gradient-to-br from-primary/20 to-primary/5 rounded-lg flex items-center justify-center">
              <FileText className="h-5 w-5 text-primary" />
            </div>
            <div>
              <DialogTitle className="text-xl font-semibold">
                {isEditing ? 'Editar Nota' : 'Criar Nova Nota'}
              </DialogTitle>
              <DialogDescription className="text-muted-foreground mt-1">
                {isEditing
                  ? 'Atualize o conteúdo da sua nota e faça suas alterações.'
                  : 'Compartilhe seus pensamentos, ideias ou informações importantes com os membros da sala.'}
              </DialogDescription>
            </div>
          </div>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6 pt-6">
          {error && (
            <div className="rounded-lg bg-destructive/10 border border-destructive/20 p-4 text-sm text-destructive flex items-start space-x-3">
              <div className="w-2 h-2 bg-destructive rounded-full mt-2 flex-shrink-0"></div>
              <div>{error}</div>
            </div>
          )}

          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <Label htmlFor="title" className="text-sm font-semibold">
                Título da Nota
              </Label>
              <span className="text-xs text-muted-foreground">
                {register('title').name && errors.title ? '0' : (document.getElementById('title') as HTMLInputElement)?.value?.length || 0}/100
              </span>
            </div>
            <Input
              id="title"
              placeholder="Digite um título descritivo para sua nota..."
              {...register('title')}
              className={`h-12 text-base ${errors.title ? 'border-destructive focus-visible:ring-destructive/20' : 'focus-visible:ring-primary/20'}`}
            />
            {errors.title && (
              <p className="text-sm text-destructive flex items-center space-x-1">
                <span>•</span>
                <span>{errors.title.message}</span>
              </p>
            )}
          </div>

          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <Label htmlFor="content" className="text-sm font-semibold">
                Conteúdo da Nota
              </Label>
              <span className="text-xs text-muted-foreground">
                {register('content').name && errors.content ? '0' : (document.getElementById('content') as HTMLTextAreaElement)?.value?.length || 0}/5000
              </span>
            </div>
            <div className="relative">
              <textarea
                id="content"
                placeholder="Escreva o conteúdo da sua nota aqui. Você pode incluir informações detalhadas, links ou outros detalhes relevantes..."
                rows={12}
                className={`w-full rounded-lg border bg-background px-4 py-3 text-sm leading-relaxed resize-none focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 transition-all ${
                  errors.content 
                    ? 'border-destructive focus-visible:ring-destructive/20' 
                    : 'border-input focus-visible:ring-primary/20 focus-visible:border-primary/30'
                }`}
                {...register('content')}
              />
            </div>
            {errors.content && (
              <p className="text-sm text-destructive flex items-center space-x-1">
                <span>•</span>
                <span>{errors.content.message}</span>
              </p>
            )}
          </div>

          <div className="flex flex-col sm:flex-row justify-end space-y-2 sm:space-y-0 sm:space-x-3 pt-6 border-t">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isLoading}
              className="sm:w-auto w-full"
            >
              Cancelar
            </Button>
            <Button 
              type="submit" 
              disabled={isLoading}
              className="sm:w-auto w-full bg-gradient-to-r from-primary to-primary/80 hover:from-primary/90 hover:to-primary/70"
            >
              {isLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  {isEditing ? 'Atualizando...' : 'Criando...'}
                </>
              ) : (
                <>
                  <FileText className="mr-2 h-4 w-4" />
                  {isEditing ? 'Atualizar Nota' : 'Criar Nota'}
                </>
              )}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
};