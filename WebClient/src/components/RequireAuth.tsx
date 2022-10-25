import { Navigate, useLocation } from "react-router-dom"
import { useAuth } from "../services/useAuth";

type Props = {
    children: JSX.Element,
  };

const RequireAuth = ({children}: Props) => {
    const location = useLocation();
    const user = useAuth();

    if(!user){
        return <Navigate to='/login' state={{from: location}}/>
    }

    return children;
}

export {RequireAuth}