'use client'
//import { Signin } from "@/actions"
import { Input } from "@/components/ui/input"
import useAuthStore from "@/store/auth" 
import {User} from "@/store/auth"
import { useForm, SubmitHandler } from "react-hook-form"
//import { Login } from '@/components/authen/login';
import { Separator } from "@/components/ui/separator"
import { useRouter } from 'next/router';
//import Router from 'next/router'
//import { useEffect } from "react"



 





export default function Login() {
 
  const {
    register,
    handleSubmit
  } = useForm<User>()
 
  const { Signin, isLoggedIn } = useAuthStore();
  //const {login} = useAuthStore()
 // const router = useRouter();

  // useEffect(() => {
  //   // ตรวจสอบสถานะการล็อกอินเพื่อ redirect
  //   if (state.isLoggedIn) {
  //     router.push('/dashboard');
  //   }
  // }, [state.isLoggedIn]); // ทำงานเมื่อ isLoggedIn เปลี่ยนแปลง


  const onSubmit: SubmitHandler<User> = async (data:User) => {
    await Signin(data); // เรียกใช้ฟังก์ชัน Signin ใน authStore
    location.replace("/")
    // if (isLoggedIn) {
    //   router.push('/dashboard'); // redirect หลัง login สำเร็จ
    // }
    //  login(data).then((response)=>{
    
    //   if(response.Status){
    //     document.cookie = "isLoggedIn=true; path=/"; // อัพเดต cookie
    //     state.accessToken = response.token
    //     //navigate("dashboard")
    //    // router.push('/dashboard');
      
    //   }else {
    //     state.isLoggedIn = false
    //     document.cookie = "isLoggedIn=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC;"; // ลบ cookie

    //     }

    //  })
       //console.log('logged in')
  //  login(data)
 
  }

  const redirect = ()=>{
    location.replace("/register")
}

  // const handleLogOut() => {
  //     console.log('logged out')
  //     logout()
  //     }
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