/*
 * This file contains functions related to creating + freeing
 * objects on behalf of the go runtime
 */

#include "linear.h"
#include <stdlib.h>

extern "C" {

/* NOTE: the Golang versions of the structures must call the corresponding
 * Free functions via runtime.SetFinalize */
/* CreateCProblem allocates a new struct problem outside of Golang's
 * garbage collection. */
struct problem *CreateCProblem() {
    auto ret = new problem();
    *ret = {}; // < Clear all fields
    return ret;
}

/* CreateCModel allocates a new struct model outside of Golang's 
 * garbage collection. */
struct model *CreateCModel() {
    auto ret = new model();
    *ret = {}; // < Clear all fields
    return ret;
}

/* CreateCParameter allocates a new struct parameter outside of
 * Golang's garbage collection.*/
struct parameter *CreateCParameter() {
    return reinterpret_cast<struct parameter*>(calloc(1, sizeof(struct parameter)));
}

/* Free's a previously allocated problem and all its data */
void FreeCProblem(struct problem *p) {
    if (p->y != nullptr) {
        free(p->y);
        p->y = nullptr;
    }
    if (p->x != nullptr) {
        free(p->x);
        p->x = nullptr;
    }
    delete p;
}

/* free's a model with libsvm's internal routines */
void FreeCModel(struct model *m) {
    free_model_content(m);
    delete m;
}

/* free's a parameter via libsvm */
void FreeCParameter(struct parameter *p) {
    if (p == nullptr) {
        return;
    }
    free(p);
}

/* Allocates a vector of doubles for storing target values
 * outside of Go's garbage collection */
int AllocateLabelsForProblem (struct problem *p, int numValues) {
    p->y = reinterpret_cast<double *>(malloc(sizeof(double) * numValues));
    return p->y == nullptr;
}

/* Utility method used to set the target value for a particular
 * input row */
void AssignLabelForProblem(struct problem *p, int i, double d) {
    p->y[i] = d;
}

/* Allocates a buffer of input rows and inserts the per-row values */
int RiffleFeatures(struct problem *p, int num_offsets, int* row_offsets, struct feature_node *features) {

    // Allocate space for the feature node buffer.
    p->x = reinterpret_cast<struct feature_node**>(
        calloc(num_offsets, sizeof(struct feature_node *))
    );
    if (p->x == nullptr) {
        return -1;
    }

    for (int i = 0; i < num_offsets; i++) {
        int offset = row_offsets[i];
        p->x[i] = features + offset;
    }
    return 0;
}

} /* extern "C" */
