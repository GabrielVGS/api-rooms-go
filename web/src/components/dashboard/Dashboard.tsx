import React, { useState, useEffect } from 'react';
import { Users, Plus } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { notesApi, roomApi, usersApi } from '@/services/api';
import type { Room, Reservation } from '@/types/api';


export const Dashboard: React.FC = () => {
  const [stats, setStats] = useState({
    totalRooms: 0,
    totalUsers: 0
  });
  const [rooms, setRooms] = useState<Room[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const roomsData = await roomApi.getRooms()
        const usersData = await usersApi.getUsers()
        console.log(usersData)




        setRooms(roomsData);
        setStats({
          totalRooms: roomsData.length,
          totalUsers: usersData.length
        });
      } catch (error) {
        console.error('Failed to fetch dashboard data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);


  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="text-muted-foreground">Carregando...</div>
      </div>
    );
  }

  console.log(stats)

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="text-muted-foreground">
          Bem vindo de volta
        </p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Salas totais</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.totalRooms}</div>
            <p className="text-xs text-muted-foreground">
              Salas disponíveis
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Usários totais</CardTitle>
            <Users className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.totalUsers}</div>
            <p className="text-xs text-muted-foreground">
              Usuários do sistema
            </p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Ações rápidas</CardTitle>
            <CardDescription>Ações</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Link to="/reservations">
              <Button className="w-full justify-start">
                <Plus className="mr-2 h-4 w-4" />
                Nova reserva
              </Button>
            </Link>
            <Link to="/rooms">
              <Button variant="outline" className="w-full justify-start">
                <Users className="mr-2 h-4 w-4" />
                Gerenciar salas
              </Button>
            </Link>
          </CardContent>
        </Card>


      </div>
    </div>
  );
};