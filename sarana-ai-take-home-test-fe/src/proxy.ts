import { NextRequest, NextResponse } from "next/server";

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

function getTokenFromRequest(request: NextRequest): string | null {
  const cookieToken = request.cookies.get('sarana-notes-app')?.value;
  if (cookieToken) return cookieToken;

  const authHeader = request.headers.get('authorization');
  if (authHeader?.startsWith('Bearer ')) {
    return authHeader.substring(7);
  }

  return null;
}

async function validateToken(token: string): Promise<boolean> {
  try {
    const response = await fetch(`${API_URL}/notes?limit=1`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      },
      cache: 'no-store',
    });

    return response.status === 200;
  } catch (error) {
    console.error('Token validation error:', error);
    return false;
  }
}

export async function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;

  const protectedRoutes = ['/dashboard'];
  const isProtectedRoute = protectedRoutes.some(route => pathname.startsWith(route));

  const authRoutes = ['/login', '/register'];
  const isAuthRoute = authRoutes.includes(pathname);

  // let token: string | null = null;

  // const persistedCookie = request.cookies.get('sarana-notes-app')?.value;
  // if (persistedCookie) {
  //   const parsed = JSON.parse(persistedCookie);
  //   token = parsed.state?.token || null;
  // }

  const token = getTokenFromRequest(request);

  if (!token && isProtectedRoute) {
    const url = request.nextUrl.clone();
    url.pathname = '/login';
    url.searchParams.set('redirect', pathname);
    return NextResponse.redirect(url);
  }

  if (token && isProtectedRoute) {
    const isValid = await validateToken(token);

    if (!isValid) {
      const url = request.nextUrl.clone();
      url.pathname = '/login';
      url.searchParams.set('redirect', pathname);
      url.searchParams.set('reason', 'session_expired');

      const response = NextResponse.redirect(url);
      response.cookies.delete('sarana-notes-app');
      return response;
    }
  }

  if (token && isAuthRoute) {
    const isValid = await validateToken(token);
    if (isValid) {
      const url = request.nextUrl.clone();
      url.pathname = '/dashboard';
      return NextResponse.redirect(url);
    }
  }

  const response = NextResponse.next();

  return response;
}

export const config = {
  matcher: [
    /*
     * Match all request paths except:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     */
    '/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
  ],
};