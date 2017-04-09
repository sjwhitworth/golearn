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
    auto ret = new parameter();
    *ret = {};
    return ret;
}

/* Free's a previously allocated problem and all its data */
void FreeCProblem(struct problem *p) {
    if (p->y != nullptr) {
        free(p->y);
        p->y = nullptr;
    }
    if (p->x != nullptr) {
        // l is the total count of rows in the problem
        // n is the number of values in each row
        for (int i = 0; i < p->l; i++) {
            if (p->x[i] != nullptr) {
                free(p->x[i]);
                p->x[i] = nullptr;
            }
        }
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
    destroy_param(p);
    delete p;
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

/* Returns a feature node for a particular row and column. */
struct feature_node *GetFeatureNodeForIndex(struct problem *p, int i, int j) {
    return &(p->x[i][j]);
}

/* Allocates a buffer of input rows and the values to fill them. */
int AllocateFeatureNodesForProblem(struct problem *p, 
        int numSamples, int numValues) {

    numValues++; // Extend for terminating element
    p->x = reinterpret_cast<struct feature_node **>(
        calloc(numSamples, sizeof(struct feature_node *))
    );
    if (p->x == nullptr) {
        return -1;
    }

    for (int i = 0; i < numSamples; i++) {
        p->x[i] = reinterpret_cast<struct feature_node *>(
            calloc(numValues, sizeof(struct feature_node))
        );
        if (p->x[i] == nullptr) {
            return -1;
        }
        // Write the special terminating element, which signals
        // to libsvm that there's no more data available on this row.
        p->x[i][numValues-1].index = -1;
    }
    return 0;
}

} /* extern "C" */
