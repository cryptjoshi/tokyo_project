'use client'
import React from "react"
import { useRouter } from 'next/router';
//import { Signin } from "@/actions"
import { Input } from "@/components/ui/input"
import useAuthStore from "@/store/auth" 
import {User} from "@/store/auth"
import { useForm, SubmitHandler } from "react-hook-form"
//import { Login } from '@/components/authen/login';
import { Separator } from "@/components/ui/separator"
//import { useRouter } from 'next/router';
//import Router from 'next/router'
//import { useEffect } from "react"
import { ToastAction } from "@/components/ui/toast"

import { useToast } from "@/hooks/use-toast"

 





export default function Login() {
  const router = useRouter();
  const { Signin, isLoggedIn } = useAuthStore();
 
// const useStore = <T, F>(
//   store: (callback: (state: T) => unknown) => unknown,
//   callback: (state: T) => F
// ) => {
//   const result = store(callback) as F;
//   const [data, setData] = React.useState<F>();

//   React.useEffect(() => {
//     setData(result);
//   }, [result]);

//   return data;
// };
   
  const {toast} = useToast()
  const {
    register,
    handleSubmit
  } = useForm<User>()
 
  
  //const {login} = useAuthStore()
 

  // useEffect(() => {
  //   // ตรวจสอบสถานะการล็อกอินเพื่อ redirect
  //   if (state.isLoggedIn) {
  //     router.push('/dashboard');
  //   }
  // }, [state.isLoggedIn]); // ทำงานเมื่อ isLoggedIn เปลี่ยนแปลง
  //const store = useStore(useAuthStore, (state:AuthStore) => state)
 
  const onSubmit: SubmitHandler<User> = async (data:User) => {
    try {

    
     const response = await Signin(data); // เรียกใช้ฟังก์ชัน Signin ใน authStore
    
     
     if (response) {
  
        router.push('/dashboard'); // เปลี่ยนไปที่หน้า dashboard เมื่อเข้าสู่ระบบสำเร็จ
    } else {
  
      toast({
        variant: "destructive",
        title: "มีข้อผิดพลาด",
        description: "เข้าสู่ระบบผิดพลาด!",
        action: <ToastAction altText="ผิดพลาด">{"ผิดพลาด"}</ToastAction>,
      })
    }
    }
    catch
    {}
  }

  const redirect = ()=>{
    location.replace("/register")
}

  React.useEffect(() => {
  if (isLoggedIn) {
      router.push('/dashboard'); // หากผู้ใช้เข้าสู่ระบบอยู่แล้ว นำทางไปยังหน้า dashboard
  }
}, [isLoggedIn, router]);

  return (
    <>
    <form onSubmit={handleSubmit(onSubmit)}>
    <div className="bg-gray-100 min-h-screen flex items-center justify-center p-6">
      <div className="bg-white shadow-lg rounded-lg max-w-md mx-auto">
        <div className="px-6 py-4">
          <h2 className="text-gray-700 text-3xl font-semibold">Login</h2>
          <p className="mt-1 text-gray-600">Please login to your account.</p>
        </div>
        <div className="px-6 py-4">
            <div className="mt-4">
              <label className="block text-gray-700" htmlFor="username">
                Username
              </label>
              <Input
                type="text"
                id="username"
                className="mt-2 rounded w-full px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300"
                placeholder=""
                required
                defaultValue=""  {...register("username", { required: true })} 
              />
            </div>
            <div className="mt-4">
              <label className="block text-gray-700" htmlFor="password">
                Password
              </label>
              <Input
                type="password"
                id="password"
                className="mt-2 rounded w-full px-3 py-2 text-gray-700 bg-gray-200 outline-none focus:bg-gray-300"
                required
                defaultValue="" {...register("password", { required: true })} 
              />
            </div>
            <div className="mt-6">
              <button type="submit" className="py-2 px-4 bg-gray-700 text-white rounded hover:bg-gray-600 w-full">
                Login
              </button>
            </div>
            <Separator className="my-4" />
            <div className="mt-3">
            <button onClick={redirect} className="py-2 px-4 bg-gray-700 text-white rounded hover:bg-gray-600 w-full">
                Register
              </button>
            </div>
        </div>
      </div>
    </div>
    </form>
    </>
  )
}