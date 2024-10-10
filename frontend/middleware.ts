import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const url = request.nextUrl.clone();

  // ดึงค่าจาก cookies ที่บันทึกสถานะล็อกอิน
  const isLoggedIn = request.cookies.get('isLoggedIn');

  if (!isLoggedIn) {
    url.pathname = '/login';
    return NextResponse.redirect(url);
  }
  if (url.pathname === '/') {
    if (!isLoggedIn) {
      url.pathname = '/login';
    } else {
      url.pathname = '/dashboard';
    }
    return NextResponse.redirect(url);
  }  
  return NextResponse.next();
}

export const config = {
  matcher: ['/', '/dashboard'],
};
