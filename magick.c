#include <errno.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>

#include <magick/api.h>

static void magick_fix_signal(int signum)
{
    struct sigaction st;
    
    if (sigaction(signum, NULL, &st) < 0) {
        goto fix_signal_error;
    }
    st.sa_flags |= SA_ONSTACK;
    if (sigaction(signum, &st,  NULL) < 0) {
        goto fix_signal_error;
    }
    return;
fix_signal_error:
        fprintf(stderr, "error fixing handler for signal %d, please "
                "report this issue to "
                "https://github.com/rainycape/magick: %s\n",
                signum, strerror(errno));
}

static void magick_fix_signals()
{
#if defined(SIGCHLD)
    magick_fix_signal(SIGCHLD);
#endif
#if defined(SIGHUP)
    magick_fix_signal(SIGHUP);
#endif
#if defined(SIGINT)
    magick_fix_signal(SIGINT);
#endif
#if defined(SIGQUIT)
    magick_fix_signal(SIGQUIT);
#endif
#if defined(SIGABRT)
    magick_fix_signal(SIGABRT);
#endif
#if defined(SIGFPE)
    magick_fix_signal(SIGFPE);
#endif
#if defined(SIGTERM)
    magick_fix_signal(SIGTERM);
#endif
#if defined(SIGBUS)
    magick_fix_signal(SIGBUS);
#endif
#if defined(SIGSEGV)
    magick_fix_signal(SIGSEGV);
#endif
#if defined(SIGXCPU)
    magick_fix_signal(SIGXCPU);
#endif
#if defined(SIGXFSZ)
    magick_fix_signal(SIGXFSZ);
#endif
}

void magick_init()
{
#ifdef _MAGICK_USES_GM
    InitializeMagick(NULL);
#else
    MagickCoreGenesis(NULL, MagickTrue);
#endif

    magick_fix_signals();
}

void magick_cleanup()
{
#ifdef _MAGICK_USES_GM
    DestroyMagick();
#else
    MagickCoreTerminus();
#endif
}
