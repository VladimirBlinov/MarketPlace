import { Route, Routes } from 'react-router-dom';
import './App.css';
import Home from './components/Home';
import { Layout } from './components/Layout';
import { Login } from './components/Login';
import ProductList from './components/ProductList';
import {CreateProduct} from './components/CreateProduct';
import Profile from './components/Profile';
import { Register } from './components/Register';
import { RequireAuth } from './components/RequireAuth';
import { AuthProvider } from './services/authProvider';
import { EditProduct } from './components/EditProduct';
import { Product } from './components/Product';

export default function App () {
  return (

    <AuthProvider>
        <Routes>
          <Route path={"/"} element={<Layout />} >
            
            <Route index element={<Home />} />
            <Route path="login" element={<Login />} />
            <Route path="register" element={<Register/>} />
            <Route path="profile" element={
              <RequireAuth>
                <Profile />
              </RequireAuth>
            } />
            <Route path="products" element={
              <RequireAuth>
                <ProductList />
              </RequireAuth>
            } />
            <Route path="products/edit/:product_id" element={
              <RequireAuth>
                <EditProduct />
              </RequireAuth>
            } />
            <Route path="products/:product_id" element={
              <RequireAuth>
                <Product />
              </RequireAuth>
            } />
            <Route path="products/new" element={
              <RequireAuth>
                <CreateProduct />
              </RequireAuth>
            } />
          </Route>
        </Routes>
    </AuthProvider>
  );
};

