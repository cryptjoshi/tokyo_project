 
 
//import useAuthStore from "@/store/auth"

 

type User = {
    username: string;
    password: string;
}


export const Signin = async (body:User) =>{
     
    try {
       // const state = useAuthStore()

        const response = await fetch("https://backend.paribrand.shop/api/v1/users/login", { method: 'POST',
            headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            //'Authorization': 'Bearer ' +  token
            },
            body: JSON.stringify({"username":body.username,password:body.password})
          })
        const data = await response.json()
        //set({ isLoggedIn: true });
        //set({ products: data.products, isLoading: false })
        //const userLocalStorage = localStorage.getItem('accessToken');
        //if (userLocalStorage) {
        if (data.status){
           // set({ isLoggedIn: true });
         //  state.setAccessToken(data.token)
         //  state.isLoggedIn = data.status
        }
    }
    catch (err:unknown) {
        console.error(err);
    }
}