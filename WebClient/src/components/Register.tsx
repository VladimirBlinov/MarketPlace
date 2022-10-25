import { useForm } from "react-hook-form";
import './Register.scss'
import './Forms.scss'
import { logon } from "../services/auth.service";
import { useNavigate  } from "react-router-dom";
import { useState } from "react";

export const Register = () => {
    const [message, setMessage] = useState<string>("");
    
    const {
        register,
        formState:{
            errors,
        },
        watch,
        handleSubmit,
    } = useForm({
        mode: "onBlur"
    });
    
    let navigate = useNavigate();
    
    const onSubmit = (data: any) => {
        logon(data.email, data.password).then(
            (response) => {
                setMessage(response.data.message);
                navigate("/login", {replace: true});
            })
    }

    return (
        <div className="form-plug">
        <section className="form-register">
        <h2>Sign up</h2>
        <form onSubmit={handleSubmit(onSubmit)}>
            <label className="register">
                Email:
                <input type="text" placeholder="Email" {...register("email", 
                    {required: "Обязательное поле", 
                    pattern: {value:/^\S+@\S+$/i,
                    message: "Неправильный Email"}
                    })} />
            </label>
            <div className="errmsg">
                {errors?.email && errors?.email?.message?.toString()}
            </div>

            <label className="register">
                Password:
                <input type="password" placeholder="Password" {...register("password", {required: "Обязательное поле", min: 8})} />
            </label>
            <div className="errmsg">
                {errors?.password && errors?.password?.message?.toString()}
            </div>

            <label className="register">
                Confirm Password
                <input type="password" placeholder="Confirm Password" {...register("confirm_password",
                 {required: true,
                  validate: (val: string) => {
                    if (watch('password') !== val) {
                        return "Пароль не совпадает";
                    }}})
                } />
            </label>
            <div className="errmsg">
                {errors?.confirm_password && errors?.confirm_password?.message?.toString()}
            </div>
            
            <input type="submit" value="Submit" />
        </form>
       </section>
       </div>
    )
}