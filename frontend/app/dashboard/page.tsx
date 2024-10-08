'use client'
import { Button } from "@/components/ui/button";
import useAuthStore from "@/store/auth"

export default function Dashboar(){

    const {logout} = useAuthStore()
    return (
        <>
       
        <div className="mt-6 p-5">
        Dashboard
       
        </div>
        <div className="mt-6 p-5">
        <Button onClick={()=>{logout()}}>Logout</Button>
        </div>

        </>
    )
}