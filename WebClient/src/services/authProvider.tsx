import { createContext, useState } from "react";
import IUser from "../entities/User";

type Props = {
    children: JSX.Element,
  };

export const AuthContext = createContext<any>(null);

export const AuthProvider = ({children}: Props) => {
    const [user, setUser] = useState<IUser | undefined>(undefined);

    const signin = (newUser: IUser, cb: () => void) =>{
        setUser(newUser);
        cb();
    }

    const signout = (cb: () => void) => {
        setUser(undefined);
        cb();
    }

    const values = {user, signin, signout}

    return (<AuthContext.Provider value ={values}>
        {children}
    </AuthContext.Provider>)
}