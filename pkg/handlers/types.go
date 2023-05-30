// SPDX-FileCopyrightText: 2023 Steffen Vogel <post@steffenvogel.de>
// SPDX-License-Identifier: Apache-2.0

package handlers

type part struct {
	Number int64  `json:"number"`
	ETag   string `json:"etag"`
	URL    string `json:"url,omitempty"`
	Length int    `json:"length,omitempty"`
	Offset uint64 `json:"offset,omitempty"`
}
