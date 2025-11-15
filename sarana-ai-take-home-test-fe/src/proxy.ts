import { NextRequest, NextResponse } from "next/server";
import { getAuthCookie } from "./lib/cookie";

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const protectedRoutes = [
  '/dashboard'
];

const authRoutes = [
  '/login',
  '/register'
];

async function validateToken(): Promise<boolean> {
  try {
    const token = getAuthCookie()?.token;
    const res = await fetch(`${API_URL}/notes`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
      },
    });
    return res.ok;
  } catch {
    return false;
  }
}

export async function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;

  const isProtectedRoute = protectedRoutes.some(route => pathname.startsWith(route));
  const isAuthRoute = authRoutes.includes(pathname);

  if (isProtectedRoute && !(await validateToken())) {
    return NextResponse.next();
  }

  if (isAuthRoute) {
    return NextResponse.next();
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
  ],
};
