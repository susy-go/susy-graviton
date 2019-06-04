/*
  This file is part of sofash.

  sofash is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  sofash is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MSRCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with sofash.  If not, see <http://www.gnu.org/licenses/>.
*/

/** @file sofash.h
* @date 2015
*/
#pragma once

#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <stddef.h>
#include "compiler.h"

#define SOFASH_REVISION 23
#define SOFASH_DATASET_BYTES_INIT 1073741824U // 2**30
#define SOFASH_DATASET_BYTES_GROWTH 8388608U  // 2**23
#define SOFASH_CACHE_BYTES_INIT 1073741824U // 2**24
#define SOFASH_CACHE_BYTES_GROWTH 131072U  // 2**17
#define SOFASH_EPOCH_LENGTH 30000U
#define SOFASH_MIX_BYTES 128
#define SOFASH_HASH_BYTES 64
#define SOFASH_DATASET_PARENTS 256
#define SOFASH_CACHE_ROUNDS 3
#define SOFASH_ACCESSES 64
#define SOFASH_DAG_MAGIC_NUM_SIZE 8
#define SOFASH_DAG_MAGIC_NUM 0xFEE1DEADBADDCAFE

#ifdef __cplusplus
extern "C" {
#endif

/// Type of a seedhash/blockhash e.t.c.
typedef struct sofash_h256 { uint8_t b[32]; } sofash_h256_t;

// convenience macro to statically initialize an h256_t
// usage:
// sofash_h256_t a = sofash_h256_static_init(1, 2, 3, ... )
// have to provide all 32 values. If you don't provide all the rest
// will simply be unitialized (not guranteed to be 0)
#define sofash_h256_static_init(...)			\
	{ {__VA_ARGS__} }

struct sofash_light;
typedef struct sofash_light* sofash_light_t;
struct sofash_full;
typedef struct sofash_full* sofash_full_t;
typedef int(*sofash_callback_t)(unsigned);

typedef struct sofash_return_value {
	sofash_h256_t result;
	sofash_h256_t mix_hash;
	bool success;
} sofash_return_value_t;

/**
 * Allocate and initialize a new sofash_light handler
 *
 * @param block_number   The block number for which to create the handler
 * @return               Newly allocated sofash_light handler or NULL in case of
 *                       ERRNOMEM or invalid parameters used for @ref sofash_compute_cache_nodes()
 */
sofash_light_t sofash_light_new(uint64_t block_number);
/**
 * Frees a previously allocated sofash_light handler
 * @param light        The light handler to free
 */
void sofash_light_delete(sofash_light_t light);
/**
 * Calculate the light client data
 *
 * @param light          The light client handler
 * @param header_hash    The header hash to pack into the mix
 * @param nonce          The nonce to pack into the mix
 * @return               an object of sofash_return_value_t holding the return values
 */
sofash_return_value_t sofash_light_compute(
	sofash_light_t light,
	sofash_h256_t const header_hash,
	uint64_t nonce
);

/**
 * Allocate and initialize a new sofash_full handler
 *
 * @param light         The light handler containing the cache.
 * @param callback      A callback function with signature of @ref sofash_callback_t
 *                      It accepts an unsigned with which a progress of DAG calculation
 *                      can be displayed. If all goes well the callback should return 0.
 *                      If a non-zero value is returned then DAG generation will stop.
 *                      Be advised. A progress value of 100 means that DAG creation is
 *                      almost complete and that this function will soon return succesfully.
 *                      It does not mean that the function has already had a succesfull return.
 * @return              Newly allocated sofash_full handler or NULL in case of
 *                      ERRNOMEM or invalid parameters used for @ref sofash_compute_full_data()
 */
sofash_full_t sofash_full_new(sofash_light_t light, sofash_callback_t callback);

/**
 * Frees a previously allocated sofash_full handler
 * @param full    The light handler to free
 */
void sofash_full_delete(sofash_full_t full);
/**
 * Calculate the full client data
 *
 * @param full           The full client handler
 * @param header_hash    The header hash to pack into the mix
 * @param nonce          The nonce to pack into the mix
 * @return               An object of sofash_return_value to hold the return value
 */
sofash_return_value_t sofash_full_compute(
	sofash_full_t full,
	sofash_h256_t const header_hash,
	uint64_t nonce
);
/**
 * Get a pointer to the full DAG data
 */
void const* sofash_full_dag(sofash_full_t full);
/**
 * Get the size of the DAG data
 */
uint64_t sofash_full_dag_size(sofash_full_t full);

/**
 * Calculate the seedhash for a given block number
 */
sofash_h256_t sofash_get_seedhash(uint64_t block_number);

#ifdef __cplusplus
}
#endif
