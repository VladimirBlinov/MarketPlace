import {Link, Outlet, useNavigate} from 'react-router-dom'
import './Header.scss'
import { useAuth } from '../services/useAuth';

const Layout = () => {
    const {user} = useAuth();
    const navigate = useNavigate();
    const {signout} = useAuth();

    return(
        <>
        <header className="page-header">
        <h1 className='page-header-h1'>MarketPlace</h1>
        </header>
        <menu>
            <Link to="/">Home</Link>
            
            {user ? (
            <>
            <Link to="/profile">Profile</Link>
            <Link to="/products">Products</Link>
            <Link to="/login" onClick={() => signout(() => navigate('/login', {replace: true}))}>Logout</Link>
            </>
            ) : (
            <>
            <Link to="/register">Register</Link>
            <Link to="/login">Login</Link>
            </>
            )
            }
        </menu>
        <Outlet />
        </>
    )
}

export {Layout}