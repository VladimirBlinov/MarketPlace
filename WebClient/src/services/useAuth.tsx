import { useContext } from "react";
import { AuthContext } from "./authProvider";


export function useAuth(){
    return useContext(AuthContext)
}