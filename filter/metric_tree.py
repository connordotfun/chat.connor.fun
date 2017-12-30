import numpy as np

class Tree:

    default_tolerance = 3
    root = None

    class Node(object):
        def __init__(self, word):
            self.word = word
            self.uses = 0
            self.children = {}

    def add_word(self, word):
        if self.root == None:
            self.root = self.Node(word)
        else:
            self.__add_word_rec(self.Node(word))

    def __add_word_rec(self, new_node):
        parent = self.root

        dist = self.__distance(new_node.word, parent.word)

        while dist in parent.children:
            parent = parent.children[dist]
            dist = self.__distance(new_node.word, parent.word)
        
        parent.children[dist] = new_node

    def get_closest(self, word, tolerance=default_tolerance):
        results = []
        corrections = self.__get_closest_rec(word, self.root, tolerance)
        for node in corrections:
            results.append(node.word)
        dist_list = [self.__distance(results[i], word) for i in range(len(results))]

        sorted_results = [word for _,word in sorted(zip(dist_list, results))]
        return sorted_results[0]

    def __get_closest_rec(self, word, root, tolerance):
        dist = self.__distance(word, root.word)
        similar_words = []
        if dist < tolerance:
            similar_words.append(root)

        min_val = max(1, dist - tolerance)
        max_val = dist + tolerance

        for i in range(min_val, max_val+1):
            if i in root.children:
                similar_words += self.__get_closest_rec(word, root.children[i], tolerance)

        return similar_words

    def __distance(self, s1, s2):
        if (len(s1) == 0):
            return len(s2)
        if (len(s2) == 0):
            return len(s1)
        dist = np.zeros(shape=(len(s1)+1, len(s2)+1))

        for i in range(1, len(s1)+1):
            dist[i][0] = i

            for j in range(1, len(s2)+1):
                cost = not s1[i - 1] == s2[j - 1]
                if (i == 1):
                    dist[0][j] = j

                dist[i][j] = self.self_min( \
                    dist[i - 1][j] + 1, \
                    [dist[i][j - 1] + 1, \
                    dist[i - 1][j - 1] + cost])
                if (i > 1 and j > 1 and s1[i - 1] == s2[j - 2] and s1[i - 2] == s2[j - 1]):
                    dist[i][j] = min(dist[i][j], dist[i - 2][j - 2] + cost)
        
        return int(dist[len(s1)][len(s2)])

    def self_min(self, s1, vals):
        if len(vals) == 0:
            return s1
        current_min = s1
        for i in vals:
            current_min = min(current_min, i)
        return current_min

# tree = Tree()

# # tree.add_word("test")
# tree.add_word("testt")
# tree.add_word("tst")

# tree.add_word("Test")
# tree.add_word("tesT")


# # print(tree.root.children[1].word)
# print(tree.get_closest("test"))


# # print(tree.default_tolerance)