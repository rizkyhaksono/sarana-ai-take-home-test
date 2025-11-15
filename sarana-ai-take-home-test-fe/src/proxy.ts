import { NextRequest, NextResponse } from "next/server";

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

const protectedRoutes = [
  '/dashboard',
  '/logs'
];

const authRoutes = [
  '/login',
  '/register'
];

function getTokenFromRequest(request: NextRequest): string | null {
  const cookieHeader = request.cookies.get('sarana-note-app-auth');
  return cookieHeader?.value || null;
}

async function validateToken(token: string): Promise<boolean> {
  try {
    const res = await fetch(`${API_URL}/me`, {
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

  // Get token from cookies
  const token = getTokenFromRequest(request);

  // Check if token exists and is valid
  const isAuthenticated = token ? await validateToken(token) : false;

  // If accessing protected route without valid token, redirect to login
  if (isProtectedRoute && !isAuthenticated) {
    const url = request.nextUrl.clone();
    url.pathname = '/login';
    return NextResponse.redirect(url);
  }

  // If accessing auth routes (login/register) with valid token, redirect to dashboard
  if (isAuthRoute && isAuthenticated) {
    const url = request.nextUrl.clone();
    url.pathname = '/dashboard';
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
  ],
};
