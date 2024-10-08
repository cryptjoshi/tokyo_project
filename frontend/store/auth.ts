import { create } from 'zustand';
import { persist } from 'zustand/middleware';
 
interface AuthStore {
    isLoggedIn: boolean;
    accessToken: string | undefined;
    accessTokenData: string | undefined;
    refreshToken: string | undefined;
            
    //login: () => void;
    //logout: () => void;
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

const useAuthStore = create(
    persist<AuthStore>(
        (set,get) => ({
            isLoggedIn: false,
            accessToken: undefined,
			accessTokenData: undefined,
			refreshToken: undefined,
            
            // login: () => {
            //     const userLocalStorage = localStorage.getItem('accessToken');
            //     if (userLocalStorage) {
            //         set({ isLoggedIn: true });
            //     }
            // },
            // logout: () => {
            //     set({ isLoggedIn: false });
            //     localStorage.clear();
            // },
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