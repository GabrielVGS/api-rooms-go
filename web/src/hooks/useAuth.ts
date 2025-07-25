import { useContext } from 'react';
import { AuthContext } from "@/contexts/AuthContext"

export const useAuth = () => {
  const context = useContext(AuthContext);
  console.log("contexto",context)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
