import { create } from 'zustand';
import { persist } from 'zustand/middleware';
 
interface AuthStore {
    isLoggedIn: boolean;
    accessToken: string | undefined;
    accessTokenData: string | undefined;
    refreshToken: string | undefined;
    Signin: (body: User) => Promise<void>;   
    //login: () => void;
    logout: () => void;
    setAccessToken: (accessToken: string | undefined) => void;
	setRefreshToken: (refreshToken: string | undefined) => void;
		// set tokens on the app start
	init: () => void;
	clearTokens: () => void;
}

// type AuthStore = {
// 	accessToken: string | undefined;
// 	accessTokenData: TokenData | undefined;
// 	refreshToken: string | undefined;

// 	actions: {
// 		setAccessToken: (accessToken: string | undefined) => void;
// 		setRefreshToken: (refreshToken: string | undefined) => void;
// 		// set tokens on the app start
// 		init: () => void;
// 		clearTokens: () => void;
// 	}
// }

export type User = {
    username:string
    password:string
  }
  


const useAuthStore = create(
    persist<AuthStore>(
        (set,get) => ({
            isLoggedIn: false,
            accessToken: undefined,
			accessTokenData: undefined,
			refreshToken: undefined,
            Signin: async (body: User) => {
                try {
                   
                  const response = await fetch("http://167.71.100.123:3003/api/v1/users/login", {
                    method: 'POST',
                    headers: {
                      'Accept': 'application/json',
                      'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ "username": body.username, "password": body.password, "prefix": "psc" }),
                  });
        
                  const data = await response.json();
                  
                  if (data.status) {
                    // อัพเดต state เมื่อ login สำเร็จ
                    set({
                      isLoggedIn: true,
                      accessToken: data.token,
                    });
                    document.cookie = "isLoggedIn=true; path=/"; // อัพเดต cookie
                    //location.replace("/dashboard")
                  } else {
                    set({ isLoggedIn: false, accessToken: "" });
                    document.cookie = "isLoggedIn=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC;"; // ลบ cookie
                  }
                } catch (error) {
                  console.error('Login failed', error);
                }
              },
           
            // login: (body:User) => {
            //     Signin(body)
            //     const userLocalStorage = localStorage.getItem('accessToken');
            //     if (userLocalStorage) {
            //         set({ isLoggedIn: true });
            //     }
            // },
            logout: () => {
                set({ isLoggedIn: false, accessToken: "" });
                document.cookie = "isLoggedIn=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC;"; // ลบ cookie
                location.replace("/")
            },
            setAccessToken: (accessToken: string | undefined) =>{
                const accessTokenData = (()=>{
                    try {
                        return accessToken ? accessToken : undefined;
                    } catch (error) {
                        console.error(error)
                        return undefined;
                    }
                })();
                set({
                    accessToken,
                    accessTokenData,
                });
            },
            setRefreshToken: (refreshToken: string | undefined) =>
                set({
                    refreshToken,
                }),
            init: () => {
                const {setAccessToken, setRefreshToken} = get();
                setAccessToken(localStorage.get('accessToken'));
                setRefreshToken(localStorage.get('refreshToken'));
    
            },
            clearTokens: () => 
                set({
                    accessToken: undefined,
                    accessTokenData: undefined,
                    refreshToken: undefined,
                }),
            
        }),
        {
            name: 'userLoginStatus',
        }
    )
);

export default useAuthStore;