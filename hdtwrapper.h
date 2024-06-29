// hdt_wrapper.h
#ifndef HDT_WRAPPER_H
#define HDT_WRAPPER_H

#ifdef __cplusplus
extern "C"
{

#endif

    int generateHDTWrapper(const char *filename, const char *baseurl, const char *outfile);
    int searchWrapper(const char *filename, const char *s, const char *p, const char *o, char *resultBuffer, int bufferSize);

#ifdef __cplusplus
}
#endif

#endif // HDT_WRAPPER_H
