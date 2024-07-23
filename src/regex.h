/**
 * Copyright 2018-2024 John Chadwick <john@jchw.io>
 *
 * Permission to use, copy, modify, and/or distribute this software for any purpose
 * with or without fee is hereby granted, provided that the above copyright notice
 * and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
 * REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
 * FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
 * INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
 * OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
 * TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
 * THIS SOFTWARE.
 */

#ifndef RE_H
#define RE_H

#include "common.h"

typedef struct _REGEX REGEX;

REGEX *ReParse(LPCSTR pattern);
BOOL ReMatch(REGEX *pattern, LPCSTR text);
LPSTR ReReplace(REGEX *find, LPCSTR replace, LPCSTR text);

int ReGetNumCaptures(REGEX *pattern);
int ReGetCaptureLen(REGEX *pattern, int i);
void ReGetCaptureData(REGEX *pattern, int i, LPSTR buffer);

#endif
